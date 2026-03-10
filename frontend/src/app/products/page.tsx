"use client"; // このファイルがブラウザで動く「クライアントコンポーネント」であることを示します

import { useEffect, useState } from "react";
import axios from "axios";

// TypeScriptの「型定義」：商品データの形を定義して、型ミスを防ぎます
// バックエンドのDBカラム（id, name, price, description, created_at）と対応しています
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

  // 【2】useEffect：画面が表示された瞬間に実行したい処理を書きます
  // [] は「最初の1回だけ実行する」という意味です
  useEffect(() => {
    const fetchProducts = async () => {
      try {
        // バックエンドの /products エンドポイントにGETリクエストを送ります
        // axios.get<Product[]>(...) で「Product型の配列が返ってくる」と教えています
        const res = await axios.get<Product[]>("http://localhost:8080/products");
        setProducts(res.data); // 取得したデータをStateに保存
      } catch {
        // 通信に失敗した場合（バックエンドが起動していない、エラーレスポンスなど）
        setError("商品データの取得に失敗しました。バックエンドが起動しているか確認してください。");
      } finally {
        // 成功・失敗どちらの場合も、読み込み中フラグをfalseにします
        setLoading(false);
      }
    };

    fetchProducts();
  }, []);

  // 【3】JSX：画面の見た目を返します
  return (
    <main className="max-w-4xl mx-auto px-4 py-8">
      {/* ヘッダー部分 */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800">商品一覧</h1>
        <p className="text-gray-500 mt-1">登録されている商品を表示しています</p>
      </div>

      {/* 状態に応じた表示の切り替え（条件付きレンダリング） */}
      {loading && <p className="text-gray-500">読み込み中...</p>}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-600 text-sm">{error}</p>
        </div>
      )}

      {!loading && !error && products.length === 0 && (
        <p className="text-gray-400">商品が登録されていません。</p>
      )}

      {/* 商品リストの表示 */}
      {/* map関数を使って、productsの中身を1つずつループして表示します */}
      <ul className="space-y-4">
        {products.map((product) => (
          <li
            key={product.id}
            className="bg-white border border-gray-200 rounded-xl p-5 shadow-sm hover:shadow-md transition"
          >
            {/* 商品名と価格を横並びに表示 */}
            <div className="flex items-start justify-between gap-4">
              <div className="flex items-center gap-2">
                {/* IDバッジ */}
                <span className="text-xs text-gray-400 bg-gray-100 px-2 py-0.5 rounded-full">
                  #{product.id}
                </span>
                <h2 className="text-xl font-semibold text-gray-800">{product.name}</h2>
              </div>
              {/* 価格は目立つように大きめ・緑色で表示 */}
              {/* toLocaleString("ja-JP") で 1000 → 1,000 のようにカンマ区切りにします */}
              <span className="text-lg font-bold text-green-600 whitespace-nowrap">
                ¥{product.price.toLocaleString("ja-JP")}
              </span>
            </div>

            {/* 商品の説明文 */}
            <p className="text-gray-600 text-sm mt-2">{product.description}</p>

            {/* 登録日時 */}
            <div className="mt-3 text-xs text-gray-400">
              {/* new Date(...).toLocaleString("ja-JP") で日付を見やすい形式に変換します */}
              登録日: {new Date(product.created_at).toLocaleString("ja-JP")}
            </div>
          </li>
        ))}
      </ul>
    </main>
  );
}
