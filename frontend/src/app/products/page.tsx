"use client"; // このファイルがブラウザで動く「クライアントコンポーネント」であることを示します

import { useEffect, useState } from "react";
import axios from "axios";

// TypeScriptの「型定義」：商品データの形を定義して、間違いを防ぎます
type Product = {
  id: number;
  name: string;
  price: number;
  description: string;
  created_at: string;
};

export default function ProductsPage() {
  // 【1】State（状態）の定義
  // products: バックエンドから取得した商品リストを保存します
  // loading: 読み込み中かどうかを管理します（trueなら「読み込み中...」と表示）
  // error: エラーが起きた時にそのメッセージを保存します
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 【2】データを取得する関数
  const fetchProducts = async () => {
    try {
      // axiosを使って、バックエンドのAPIに「データちょうだい」とリクエストを送ります
      const res = await axios.get<Product[]>("http://localhost:8080/products");
      // 成功したら、取得したデータをStateに保存します
      setProducts(res.data);
    } catch {
      // 失敗したら、エラーメッセージをStateに保存します
      setError("商品の取得に失敗しました。バックエンドが起動しているか確認してください。");
    } finally {
      // 成功しても失敗しても、読み込みは終わったので false にします
      setLoading(false);
    }
  };

  // 【3】useEffect：画面が表示された瞬間に実行したい処理を書きます
  useEffect(() => {
    fetchProducts();
  }, []); // [] は「最初の1回だけ実行する」という意味です

  // 【4】JSX：画面の見た目を返します
  return (
    <main className="max-w-4xl mx-auto px-4 py-8">
      {/* ヘッダー部分 */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800">商品一覧</h1>
        <p className="text-gray-500 mt-1">登録されている商品を表示しています</p>
      </div>

      {/* 状態に応じた表示の切り替え（条件付きレンダリング） */}
      {loading && <p className="text-gray-500">読み込み中...</p>}
      {error && <p className="text-red-500">{error}</p>}

      {!loading && !error && products.length === 0 && (
        <p className="text-gray-400">商品が登録されていません。</p>
      )}

      {/* 商品リストの表示 */}
      <ul className="space-y-4">
        {products.map((product) => (
          // map関数を使って、productsの中身を1つずつループして表示します
          <li
            key={product.id}
            className="bg-white border border-gray-200 rounded-xl p-5 shadow-sm hover:shadow-md transition"
          >
            {/* 商品名と価格を横並びに表示 */}
            <div className="flex items-start justify-between gap-4">
              <h2 className="text-xl font-semibold text-gray-800">{product.name}</h2>
              {/* 価格は目立つように大きめ・緑色で表示 */}
              <span className="text-lg font-bold text-green-600 whitespace-nowrap">
                ¥{product.price.toLocaleString("ja-JP")}
              </span>
            </div>

            {/* 商品の説明文 */}
            <p className="text-gray-600 text-sm mt-2">{product.description}</p>

            {/* 登録日時 */}
            <div className="mt-3 text-xs text-gray-400">
              {/* 日付を見やすい形式に変換して表示します */}
              登録日: {new Date(product.created_at).toLocaleString("ja-JP")}
            </div>
          </li>
        ))}
      </ul>
    </main>
  );
}
