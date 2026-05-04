import fs from 'fs'
import axios from 'axios'
import { GoogleGenerativeAI, GoogleGenerativeAIFetchError } from '@google/generative-ai'

const BOT_HEADER = '### 🤖 Gemini PR Review'
const apiKey = process.env.GEMINI_API_KEY1

if (!apiKey) {
  console.error('GEMINI_API_KEY1 is not set. Skipping Gemini review.')
  process.exit(0)
}

const githubToken = process.env.GITHUB_TOKEN
const repoFullName = process.env.REPO_FULL_NAME
const prNumber = process.env.PR_NUMBER

if (!githubToken || !repoFullName || !prNumber) {
  console.error('Missing required GitHub environment variables.')
  process.exit(1)
}

const github = axios.create({
  baseURL: 'https://api.github.com',
  headers: {
    Authorization: `Bearer ${githubToken}`,
    Accept: 'application/vnd.github+json',
    'Content-Type': 'application/json',
  },
})

const rulesFilePath = 'docs/REVIEW_RULES_SUMMARY.md'
const todoFilePath = 'docs/GEMINI_REVIEW_TODO.md'
let rulesSummary = ''
let reviewTodo = ''

try {
  rulesSummary = fs.readFileSync(rulesFilePath, 'utf8')
} catch (err) {
  console.warn(
    `Rules summary file "${rulesFilePath}" not found. Continuing without explicit rules summary.`
  )
}

try {
  reviewTodo = fs.readFileSync(todoFilePath, 'utf8')
} catch (err) {
  console.warn(`Review TODO file "${todoFilePath}" not found. Continuing without review TODO context.`)
}

const genAI1 = new GoogleGenerativeAI(apiKey)
const apiKey2 = process.env.GEMINI_API_KEY2
const genAI2 = apiKey2 ? new GoogleGenerativeAI(apiKey2) : null
const key1QuotaExhaustedModels = new Set()
let key1FallbackNotified = false

const MODEL_FALLBACK_LIST = [
  'gemini-2.5-flash-lite',
  'gemini-2.5-pro',
  'gemini-3.1-flash-lite-preview',
  'gemini-3-flash-preview',
  'gemini-2.5-flash',
  'gemini-2.0-flash',
]

const KEY2_MODEL = 'gemini-2.5-flash-lite'
const MAX_PATCH_CHARS = 12000
const MAX_REVIEW_FILES = 20
const SKIP_FILE_PATTERNS = [
  /\.lock$/i,
  /package-lock\.json$/i,
  /pnpm-lock\.yaml$/i,
  /yarn\.lock$/i,
  /bun\.lockb$/i,
  /go\.sum$/i,
  /vendor\//i,
  /\.min\./i,
  /(^|\/)dist\//i,
  /(^|\/)build\//i,
  /(^|\/)\.next\//i,
  /(^|\/)\.nuxt\//i,
  /(^|\/)\.output\//i,
]

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

async function notifySlack(message) {
  const webhookUrl = process.env.SLACK_WEBHOOK_URL
  if (!webhookUrl) {
    console.warn('SLACK_WEBHOOK_URL is not set. Skipping Slack notification.')
    return
  }
  try {
    await axios.post(webhookUrl, { text: message })
    console.log('Slack notification sent.')
  } catch (err) {
    console.warn('Failed to send Slack notification:', err.message)
  }
}

async function notifyKey1FallbackToKey2() {
  if (key1FallbackNotified) return
  key1FallbackNotified = true

  const exhaustedModels = [...key1QuotaExhaustedModels]
  const exhaustedModelsText =
    exhaustedModels.length > 0
      ? exhaustedModels.map((modelName) => `- ${modelName}`).join('\n')
      : '- 429 を返した KEY1 モデルは記録されていません'

  await notifySlack(
    `⚠️ GEMINI_API_KEY1 の候補モデルがすべて利用不可のため、有料の GEMINI_API_KEY2 (${KEY2_MODEL}) にフォールバックしました。\nPR: ${prNumber} / Repo: ${repoFullName}\n\nKEY1 で 429 quota exhausted を返したモデル:\n${exhaustedModelsText}`
  )
}

function shouldSkipFile(file) {
  if (!file.patch) return true
  return SKIP_FILE_PATTERNS.some((pattern) => pattern.test(file.filename))
}

function truncatePatch(patch) {
  if (patch.length <= MAX_PATCH_CHARS) return patch
  return `${patch.slice(0, MAX_PATCH_CHARS)}\n\n... [truncated]`
}

async function fetchAllPages(url) {
  const items = []
  let page = 1

  while (true) {
    const response = await github.get(url, {
      params: { per_page: 100, page },
    })

    items.push(...response.data)

    if (response.data.length < 100) break
    page += 1
  }

  return items
}

async function fetchPullRequestContext() {
  const prResponse = await github.get(`/repos/${repoFullName}/pulls/${prNumber}`)
  const files = await fetchAllPages(`/repos/${repoFullName}/pulls/${prNumber}/files`)

  return {
    pr: prResponse.data,
    files,
  }
}

function buildPrompt({ pr, file, changedFilesSummary }) {
  const prBody = pr.body?.trim() || '(本文なし)'
  const patch = truncatePatch(file.patch || '')

  return `
あなたはGitHub Actions上で動作するコードレビューボットです。

プロジェクトのレビュールール（サマリー）:
------------------------------------------------------------
${rulesSummary}
------------------------------------------------------------

未反映レビュー指摘TODO:
------------------------------------------------------------
${reviewTodo || '(TODOファイルなし、または未記載)'}
------------------------------------------------------------

PR情報:
- タイトル: ${pr.title}
- ベースブランチ: ${pr.base.ref}
- ヘッドブランチ: ${pr.head.ref}
- changed files:
${changedFilesSummary}
- PR本文:
${prBody}

あなたの役割:
- 提供されたPR情報と対象ファイルのpatchだけを根拠にレビューしてください。
- 対象ファイル以外について推測しないでください。
- 機能を壊す、またはセキュリティリスクをもたらすクリティカルな問題のみを報告してください。
- 上記のプロジェクト固有のレビュールールと制約に従ってください。
- 未反映レビュー指摘TODOに記載済みの内容は、既に人間が把握して管理しているため、同じ内容を再指摘しないでください。
- TODOに記載済みの内容でも、今回のpatchがそのリスクを明確に悪化させている場合だけ報告してください。

報告してよい問題の種類:
- bug
- security
- logic

報告してはいけないもの:
- スタイルの問題
- フォーマットや空白
- 命名規則
- パフォーマンスの微最適化
- 軽微な提案や任意の改善
- マークダウンやドキュメントの問題
- patchに含まれていない箇所への推測
- 情報不足のまま断定する指摘

重要な制約:
- 今回レビュー対象なのは以下のファイルのみです: ${file.filename}
- renameや周辺ファイルの事情は、PR情報に明示された範囲でのみ判断してください。
- patchから断定できない場合は、問題として報告しないでください。
- 問題を報告する場合、fileは必ず ${file.filename} を使ってください。

出力は必ず以下のどちらかにしてください。

1. 問題がある場合:
Issue-1:
    type: bug | security | logic
    file: ${file.filename}
    lines: (行番号または範囲)
    problem: (問題の簡潔な説明)
    reason: (なぜ問題なのか)
    suggestion: (修正案)

Issue-2:
    ...

2. 重大な問題がない場合:
重大な問題はありません。LGTM。

対象ファイル情報:
- file: ${file.filename}
- status: ${file.status}
- additions: ${file.additions}
- deletions: ${file.deletions}

===== PATCH START =====
${patch}
===== PATCH END =====
`
}

function normalizeReviewText(text) {
  return text.trim()
}

function isLgtm(text) {
  return text.includes('重大な問題はありません。LGTM。')
}

async function tryGenerate(genAI, modelName, prompt, options = {}) {
  const maxAttempts = options.maxAttempts ?? 2
  const model = genAI.getGenerativeModel({ model: modelName })
  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      console.log(`Calling Gemini model: ${modelName} (attempt ${attempt}/${maxAttempts})`)
      const result = await model.generateContent(prompt)
      console.log(`Success with model: ${modelName}`)
      return result.response.text()
    } catch (err) {
      const isFetchError = err instanceof GoogleGenerativeAIFetchError

      if (isFetchError && err.status === 429) {
        console.warn(`Model ${modelName} quota exhausted (429).`)
        if (options.recordKey1Quota) {
          key1QuotaExhaustedModels.add(modelName)
        }
        return null
      }

      if (isFetchError && err.status === 503) {
        if (attempt < maxAttempts) {
          const delay = 20000 * attempt
          console.warn(`Model ${modelName} unavailable (503). Retrying in ${delay / 1000}s...`)
          await sleep(delay)
          continue
        }
        console.warn(`Model ${modelName} still unavailable after retries.`)
        return null
      }

      console.error(`Gemini API non-retryable error (model: ${modelName}):`, err)
      throw err
    }
  }

  return null
}

async function generateReviewWithRetry(prompt) {
  for (const modelName of MODEL_FALLBACK_LIST) {
    const text = await tryGenerate(genAI1, modelName, prompt, { recordKey1Quota: true })
    if (text !== null) return text
    console.warn('Falling through to next model...')
  }

  if (genAI2) {
    console.warn(`All KEY1 models exhausted. Falling back to GEMINI_API_KEY2 (${KEY2_MODEL})...`)
    await notifyKey1FallbackToKey2()
    const text = await tryGenerate(genAI2, KEY2_MODEL, prompt)
    if (text !== null) return text
    console.error('KEY2 also failed.')
  } else {
    console.error('GEMINI_API_KEY2 is not set. Cannot fall back.')
  }

  console.error('All models and keys exhausted. Giving up for this run.')
  return null
}

function buildChangedFilesSummary(files) {
  return files
    .map((file) => `- ${file.filename} (${file.status}, +${file.additions}/-${file.deletions})`)
    .join('\n')
}

function buildFinalBody({ reviewedCount, skippedCount, results }) {
  const findings = results.filter((result) => !isLgtm(result.review))

  if (findings.length === 0) {
    return `${BOT_HEADER}

重大な問題はありません。LGTM。

- reviewed files: ${reviewedCount}
- skipped files: ${skippedCount}`
  }

  const body = findings
    .map(
      (result, index) => `#### Finding ${index + 1}: ${result.file.filename}

${result.review}`
    )
    .join('\n\n')

  return `${BOT_HEADER}

以下の重大な指摘があります。

${body}

- reviewed files: ${reviewedCount}
- skipped files: ${skippedCount}`
}

async function postComment(body) {
  await github.post(`/repos/${repoFullName}/issues/${prNumber}/comments`, { body })
  console.log('Gemini review comment posted successfully!')
}

;(async () => {
  let reviewBody

  try {
    const { pr, files } = await fetchPullRequestContext()
    const changedFilesSummary = buildChangedFilesSummary(files)

    const reviewableFiles = files.filter((file) => !shouldSkipFile(file)).slice(0, MAX_REVIEW_FILES)
    const skippedCount = files.length - reviewableFiles.length

    if (reviewableFiles.length === 0) {
      reviewBody = `${BOT_HEADER}

重大な問題はありません。LGTM。

- reviewed files: 0
- skipped files: ${files.length}`
    } else {
      const results = []

      for (const file of reviewableFiles) {
        const prompt = buildPrompt({ pr, file, changedFilesSummary })
        const text = await generateReviewWithRetry(prompt)
        results.push({
          file,
          review:
            normalizeReviewText(text || '⚠️ このファイルのレビューは一時的なAPIエラーにより取得できませんでした。'),
        })
      }

      reviewBody = buildFinalBody({
        reviewedCount: reviewableFiles.length,
        skippedCount,
        results,
      })
    }
  } catch (err) {
    console.error('Unexpected error preparing Gemini review:', err)
    process.exit(1)
  }

  try {
    await postComment(reviewBody)
  } catch (err) {
    if (err.response) {
      console.error('GitHub API error:', err.response.status, err.response.data)
    } else {
      console.error('Network error posting to GitHub:', err.message)
      console.error('Full error:', err)
    }
    process.exit(1)
  }
})()
