"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import axios from "axios";

type Post = {
  id: number;
  title: string;
  content: string;
  author: string;
  created_at: string;
};

export default function PostDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const res = await axios.get<Post>(`http://localhost:8080/posts/${params.id}`);
        setPost(res.data);
      } catch {
        setError("投稿が見つかりませんでした。");
      } finally {
        setLoading(false);
      }
    };
    fetchPost();
  }, [params.id]);

  const handleDelete = async () => {
    if (!confirm("この投稿を削除しますか？")) return;
    setDeleting(true);
    try {
      await axios.delete(`http://localhost:8080/posts/${params.id}`);
      router.push("/");
    } catch {
      alert("削除に失敗しました。");
      setDeleting(false);
    }
  };

  if (loading) return <p className="text-center mt-20 text-gray-500">読み込み中...</p>;
  if (error) return <p className="text-center mt-20 text-red-500">{error}</p>;
  if (!post) return null;

  return (
    <main className="max-w-3xl mx-auto px-4 py-8">
      <Link href="/" className="text-blue-600 hover:underline text-sm">← 一覧に戻る</Link>

      <article className="mt-6 bg-white border border-gray-200 rounded-xl p-6 shadow-sm">
        <h1 className="text-2xl font-bold text-gray-800 mb-2">{post.title}</h1>
        <div className="flex items-center gap-2 text-xs text-gray-400 mb-6">
          <span className="font-medium text-gray-500">{post.author}</span>
          <span>·</span>
          <span>{new Date(post.created_at).toLocaleString("ja-JP")}</span>
        </div>
        <p className="text-gray-700 whitespace-pre-wrap leading-7">{post.content}</p>
      </article>

      <div className="mt-4 flex justify-end">
        <button
          onClick={handleDelete}
          disabled={deleting}
          className="bg-red-500 hover:bg-red-600 disabled:opacity-50 text-white font-medium px-4 py-2 rounded-lg transition"
        >
          {deleting ? "削除中..." : "削除する"}
        </button>
      </div>
    </main>
  );
}
