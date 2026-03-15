"use client"; // このファイルがブラウザで動く「クライアントコンポーネント」であることを示します

import { useEffect, useState } from "react";
import Link from "next/link";
import axios from "axios";

// TypeScriptの「型定義」：投稿データの形を定義して、間違いを防ぎます
type Post = {
  id: number;
  title: string;
  content: string;
  author: string;
  created_at: string;
};

export default function Home() {
  // 【1】State（状態）の定義
  // posts: バックエンドから取得した投稿リストを保存します
  // loading: 読み込み中かどうかを管理します（trueなら「読み込み中...」と表示）
  // error: エラーが起きた時にそのメッセージを保存します
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 【2】データを取得する関数
  const fetchPosts = async () => {
    try {
      // axiosを使って、バックエンドのAPIに「データちょうだい」とリクエストを送ります
      const res = await axios.get<Post[]>("http://localhost:8080/posts");
      // 成功したら、取得したデータをStateに保存します
      setPosts(res.data);
    } catch {
      // 失敗したら、エラーメッセージをStateに保存します
      setError("投稿の取得に失敗しました。バックエンドが起動しているか確認してください。");
    } finally {
      // 成功しても失敗しても、読み込みは終わったので false にします
      setLoading(false);
    }
  };

  // 【3】useEffect：画面が表示された瞬間に実行したい処理を書きます
  useEffect(() => {
    fetchPosts();
  }, []); // [] は「最初の1回だけ実行する」という意味です

  // 【4】JSX：画面の見た目を返します
  return (
    <main className="max-w-3xl mx-auto px-4 py-8">
      {/* ヘッダー部分 */}
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold text-gray-800">SawakouHouse 掲示板</h1>
        <div className="flex items-center gap-3">
          <Link
            href="/map"
            className="bg-green-600 hover:bg-green-700 text-white font-medium px-4 py-2 rounded-lg transition"
          >
            マップ
          </Link>
          <Link
            href="/posts/new"
            className="bg-blue-600 hover:bg-blue-700 text-white font-medium px-4 py-2 rounded-lg transition"
          >
            新規投稿
          </Link>
        </div>
      </div>

      {/* 状態に応じた表示の切り替え（条件付きレンダリング） */}
      {loading && <p className="text-gray-500">読み込み中...</p>}
      {error && <p className="text-red-500">{error}</p>}

      {!loading && !error && posts.length === 0 && (
        <p className="text-gray-400">まだ投稿がありません。最初の投稿をしてみましょう！</p>
      )}

      {/* 投稿リストの表示 */}
      <ul className="space-y-4">
        {posts.map((post) => (
          // map関数を使って、postsの中身を1つずつループして表示します
          <li key={post.id} className="bg-white border border-gray-200 rounded-xl p-5 shadow-sm hover:shadow-md transition">
            <Link href={`/posts/${post.id}`}>
              <h2 className="text-xl font-semibold text-blue-700 hover:underline mb-1">{post.title}</h2>
            </Link>
            <p className="text-gray-600 text-sm line-clamp-2">{post.content}</p>
            <div className="mt-3 flex items-center gap-2 text-xs text-gray-400">
              <span className="font-medium text-gray-500">{post.author}</span>
              <span>·</span>
              {/* 日付を見やすい形式に変換して表示します */}
              <span>{new Date(post.created_at).toLocaleString("ja-JP")}</span>
            </div>
          </li>
        ))}
      </ul>
    </main>
  );
}
