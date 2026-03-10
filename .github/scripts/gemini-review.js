const { GoogleGenerativeAI } = require('@google/generative-ai');
const { Octokit } = require('@octokit/rest');

async function main() {
  const geminiApiKey = process.env.GEMINI_API_KEY;
  const githubToken = process.env.GITHUB_TOKEN;
  const prNumber = process.env.PR_NUMBER;
  const repoOwner = process.env.REPO_OWNER;
  const repoName = process.env.REPO_NAME;
  const modelName = process.env.GEMINI_MODEL || 'gemini-1.5-flash-latest';

  if (!geminiApiKey || !githubToken || !prNumber) {
    console.error('必要な環境変数が設定されていません');
    process.exit(1);
  }

  // GitHub APIクライアント
  const octokit = new Octokit({ auth: githubToken });

  // PR差分を取得
  const { data: files } = await octokit.pulls.listFiles({
    owner: repoOwner,
    repo: repoName,
    pull_number: prNumber,
  });

  // 差分をテキストにまとめる
  let diffText = '# プルリクエストの変更内容\n\n';
  for (const file of files) {
    diffText += `## ${file.filename}\n`;
    diffText += `変更: +${file.additions} -${file.deletions}\n`;
    if (file.patch) {
      diffText += '```diff\n' + file.patch + '\n```\n\n';
    }
  }

  // Gemini APIでレビュー
  const genAI = new GoogleGenerativeAI(geminiApiKey);
  const model = genAI.getGenerativeModel({ 
    model: modelName,
    generationConfig: {
      temperature: 0.7,
      topK: 40,
      topP: 0.95,
      maxOutputTokens: 8192,
    },
  });

  const prompt = `あなたはシニアエンジニアです。以下のプルリクエストのコードレビューを日本語で実施してください。

重要: 必ず日本語でレビューを記述してください。英語は使用しないでください。

以下の点に注目してください：
1. バグや潜在的なエラーの指摘
2. コードの可読性と保守性の向上案
3. セキュリティ上の懸念点
4. 良い実装については褒める

${diffText}`;

  const result = await model.generateContent(prompt);
  const review = result.response.text();

  // レビューコメントをPRに投稿
  await octokit.issues.createComment({
    owner: repoOwner,
    repo: repoName,
    issue_number: prNumber,
    body: `## 🤖 Gemini Code Review\n\n${review}`,
  });

  console.log('レビューコメントを投稿しました');
}

main().catch((error) => {
  console.error('エラーが発生しました:', error);
  process.exit(1);
});
