# Next.js (フロントエンド) CRUD 実装・詳細解説ガイド

このガイドでは、実際のコード例を使いながら、掲示板アプリの「投稿一覧の表示（Read）」と「新規投稿の作成（Create）」に加え、「更新（Update）」と「削除（Delete）」のフロントエンド側の仕組みを解説します。

---

## 1. Type（型定義）：データの形を保証する
`frontend/src/app/page.tsx`
TypeScript を使うことで、バックエンドから届くデータの形をあらかじめ決めておき、間違いを防ぎます。

```tsx
type Post = {
  id: number;
  title: string;
  content: string;
  author: string;
  created_at: string;
};
```

---

## 2. 【R】Read（取得）：一覧を表示する
`frontend/src/app/page.tsx`
画面が開かれた瞬間にデータを取得し、リストとして表示します。

```tsx
const [posts, setPosts] = useState<Post[]>([]);
const [loading, setLoading] = useState(true);

const fetchPosts = async () => {
  const res = await axios.get<Post[]>("http://localhost:8080/posts");
  setPosts(res.data);
  setLoading(false);
};

useEffect(() => {
  fetchPosts();
}, []);
```

---

## 3. 【C】Create（作成）：フォームから投稿する
`frontend/src/app/posts/new/page.tsx`
入力された文字を管理し、ボタンが押されたらバックエンドへ送ります。

```tsx
const [form, setForm] = useState({ title: "", content: "", author: "" });

const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();
  await axios.post("http://localhost:8080/posts", form); // 【C】
  router.push("/");
};
```

---

## 4. 【D】Delete（削除）：特定の項目を消す
`frontend/src/app/page.tsx`
削除ボタンが押されたら、バックエンドに ID を伝え、成功したら画面から消します。

```tsx
const handleDelete = async (id: number) => {
  if (!confirm("本当に削除しますか？")) return;

  try {
    // 【1】バックエンドに削除リクエストを送る
    await axios.delete(`http://localhost:8080/posts/${id}`); // 【D】
    
    // 【2】成功したら、今の posts State から削除した ID 以外を残す
    setPosts(posts.filter((post) => post.id !== id));
    alert("削除しました");
  } catch (error) {
    alert("削除に失敗しました");
  }
};

// 画面上での表示
<button onClick={() => handleDelete(post.id)}>削除</button>
```

---

## 5. 【U】Update（更新）：既存の投稿を直す
`frontend/src/app/posts/[id]/edit/page.tsx` (概念)
1.  **Fetch**: その ID の今のデータを取得して入力欄に入れる。
2.  **Edit**: 文字を書き換える (`useState`)。
3.  **Submit**: `axios.put` で変更後のデータを送る。

```tsx
const handleUpdate = async (e: React.FormEvent) => {
  e.preventDefault();
  // 【U】PUT リクエストで上書き更新
  await axios.put(`http://localhost:8080/posts/${id}`, form);
  router.push("/");
};
```

---

## 6. フロントエンド開発の重要キーワード

### axios（API 通信）
- `axios.get`: データをもらう (Read)
- `axios.post`: データを新規保存する (Create)
- `axios.put`: データを上書き更新する (Update)
- `axios.delete`: データを削除する (Delete)

### posts.filter(...)
削除後に画面を更新するために使います。「削除した ID 以外のデータを集めて、新しいリストを作る」という処理です。

---

## まとめ：フロントエンドの動き
1.  **ユーザー**：削除ボタンを押す。
2.  **handleDelete**：`axios.delete` でサーバーに命令を送る。
3.  **setPosts**：サーバーで消せたら、手元の State からも消して画面をスッキリさせる。
4.  **React**：State が変わったので、画面が自動的に再描画される。
