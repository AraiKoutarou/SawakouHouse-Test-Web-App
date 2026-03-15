"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  APIProvider,
  Map,
  AdvancedMarker,
  useMap,
  useMapsLibrary,
} from "@vis.gl/react-google-maps";
import axios from "axios";

const API_BASE = "http://localhost:8080";

type RevealData = {
  place_name: string;
  lat: number;
  lng: number;
  reason: string;
  route: { lat: number; lng: number }[];
};

// ── モックデータ（バックエンド連携前の開発用） ────────────────────
const MOCK_REVEAL: RevealData = {
  place_name: "哲学堂公園",
  lat: 35.7123,
  lng: 139.6789,
  reason:
    "あなたが「さびしい」と感じていたので、静かに自分と向き合える場所を選びました。哲学堂公園は、夜でも趣のある雰囲気があり、一人でゆっくりと考えを整理するのに最適な場所です。喧騒を離れ、自分自身と対話できる時間をお届けしたいと思いました。",
  route: [
    { lat: 35.6812, lng: 139.7671 },
    { lat: 35.6900, lng: 139.7500 },
    { lat: 35.7000, lng: 139.7200 },
    { lat: 35.7060, lng: 139.7050 },
    { lat: 35.7123, lng: 139.6789 },
  ],
};
// ─────────────────────────────────────────────────────────────────

// ルートのPolylineを描画し、地図を全体にフィットさせるコンポーネント
function RouteOverlay({ data }: { data: RevealData }) {
  const map = useMap();
  const mapsLib = useMapsLibrary("maps");

  useEffect(() => {
    if (!map || !mapsLib || data.route.length === 0) return;

    const polyline = new mapsLib.Polyline({
      path: data.route,
      geodesic: true,
      strokeColor: "#3B82F6",
      strokeOpacity: 0.85,
      strokeWeight: 5,
    });
    polyline.setMap(map);

    // ルート全体が収まるようにズームを調整
    const bounds = new window.google.maps.LatLngBounds();
    data.route.forEach((p) => bounds.extend(p));
    map.fitBounds(bounds, 80);

    return () => polyline.setMap(null);
  }, [map, mapsLib, data]);

  return null;
}

export default function ArrivedPage() {
  const { id: tripId } = useParams<{ id: string }>();
  const router = useRouter();

  const [reveal, setReveal] = useState<RevealData | null>(null);
  const [showReason, setShowReason] = useState(false);

  useEffect(() => {
    // TODO: バックエンド完成後に以下のコメントを外して切り替える
    // axios
    //   .get<RevealData>(`${API_BASE}/trips/${tripId}/reveal`)
    //   .then((r) => setReveal(r.data))
    //   .catch(console.error);

    // モック処理
    setReveal(MOCK_REVEAL);
  }, [tripId]);

  if (!reveal) {
    return (
      <div className="flex items-center justify-center h-screen bg-white">
        <p className="text-gray-500">読み込み中...</p>
      </div>
    );
  }

  return (
    <APIProvider apiKey={process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY!}>
      <div className="relative w-full h-screen overflow-hidden">
        {/* 地図（全画面） */}
        <Map
          defaultCenter={{ lat: reveal.lat, lng: reveal.lng }}
          defaultZoom={14}
          disableDefaultUI
          mapId="arrived-map"
          style={{ width: "100%", height: "100%" }}
        >
          <RouteOverlay data={reveal} />

          {/* 目的地マーカー */}
          <AdvancedMarker position={{ lat: reveal.lat, lng: reveal.lng }}>
            <div className="bg-blue-600 text-white text-sm font-bold px-3 py-1.5 rounded-full shadow-lg whitespace-nowrap">
              📍 {reveal.place_name}
            </div>
          </AdvancedMarker>
        </Map>

        {/* 下部：到着パネル */}
        <div className="absolute bottom-0 left-0 right-0 p-4">
          <div className="bg-white rounded-3xl shadow-2xl p-6">
            {/* ヘッダー */}
            <div className="text-center mb-5">
              <p className="text-2xl font-bold text-gray-900">到着しました！</p>
              <p className="text-lg text-blue-600 font-semibold mt-1">
                📍 {reveal.place_name}
              </p>
            </div>

            {!showReason ? (
              // 理由を見るボタン
              <button
                onClick={() => setShowReason(true)}
                className="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-4 rounded-2xl transition-colors text-base"
              >
                AIがここを選んだ理由を見る ✨
              </button>
            ) : (
              // 理由表示
              <div className="space-y-4">
                <div className="bg-blue-50 rounded-2xl p-4">
                  <p className="text-sm text-gray-700 leading-relaxed">
                    {reveal.reason}
                  </p>
                </div>
                <button
                  onClick={() => router.push("/")}
                  className="w-full bg-gray-900 hover:bg-gray-700 text-white font-bold py-4 rounded-2xl transition-colors text-base"
                >
                  ホームに戻る
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </APIProvider>
  );
}
