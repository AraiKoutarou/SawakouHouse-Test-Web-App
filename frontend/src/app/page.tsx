"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import axios from "axios";

type Post = {
  id: number;
  title: string;
  content: string;
  author: string;
  created_at: string;
};

export default function Home() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPosts = async () => {
    try {
      const res = await axios.get<Post[]>("http://localhost:8080/posts");
      setPosts(res.data);
    } catch {
      setError("投稿の取得に失敗しました。バックエンドが起動しているか確認してください。");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  return (
    <main className="max-w-3xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold text-gray-800">SawakouHouse 掲示板</h1>
        <Link
          href="/posts/new"
          className="bg-blue-600 hover:bg-blue-700 text-white font-medium px-4 py-2 rounded-lg transition"
        >
          新規投稿
        </Link>
      </div>

      {loading && <p className="text-gray-500">読み込み中...</p>}
      {error && <p className="text-red-500">{error}</p>}

      {!loading && !error && posts.length === 0 && (
        <p className="text-gray-400">まだ投稿がありません。最初の投稿をしてみましょう！</p>
      )}

      <ul className="space-y-4">
        {posts.map((post) => (
          <li key={post.id} className="bg-white border border-gray-200 rounded-xl p-5 shadow-sm hover:shadow-md transition">
            <Link href={`/posts/${post.id}`}>
              <h2 className="text-xl font-semibold text-blue-700 hover:underline mb-1">{post.title}</h2>
            </Link>
            <p className="text-gray-600 text-sm line-clamp-2">{post.content}</p>
            <div className="mt-3 flex items-center gap-2 text-xs text-gray-400">
              <span className="font-medium text-gray-500">{post.author}</span>
              <span>·</span>
              <span>{new Date(post.created_at).toLocaleString("ja-JP")}</span>
            </div>
          </li>
        ))}
      </ul>
    </main>
  );
}
