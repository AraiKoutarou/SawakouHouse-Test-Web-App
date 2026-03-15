"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  APIProvider,
  Map,
  AdvancedMarker,
  InfoWindow,
} from "@vis.gl/react-google-maps";
import axios from "axios";

const API_BASE = "http://localhost:8080";

type Place = {
  trip_id: string;
  place_name: string;
  lat: number;
  lng: number;
  visited_at: string;
  memory?: string;
};

// ── モックデータ（バックエンド連携前の開発用） ────────────────────
const MOCK_PLACES: Place[] = [
  {
    trip_id: "abc123",
    place_name: "哲学堂公園",
    lat: 35.7123,
    lng: 139.6789,
    visited_at: "2026-03-10T20:00:00Z",
    memory: "一人でぼーっとできた",
  },
  {
    trip_id: "def456",
    place_name: "石神井公園",
    lat: 35.7534,
    lng: 139.6011,
    visited_at: "2026-03-12T18:30:00Z",
    memory: undefined,
  },
  {
    trip_id: "ghi789",
    place_name: "神田明神",
    lat: 35.7021,
    lng: 139.7672,
    visited_at: "2026-03-14T15:00:00Z",
    memory: "静かで良かった",
  },
];
// ─────────────────────────────────────────────────────────────────

export default function HistoryPage() {
  const router = useRouter();

  const [places, setPlaces] = useState<Place[]>([]);
  const [selected, setSelected] = useState<Place | null>(null);

  useEffect(() => {
    // TODO: バックエンド完成後に以下のコメントを外して切り替える
    // axios
    //   .get<Place[]>(`${API_BASE}/places`)
    //   .then((r) => setPlaces(r.data))
    //   .catch(console.error);

    // モック処理
    setPlaces(MOCK_PLACES);
  }, []);

  return (
    <APIProvider apiKey={process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY!}>
      <div className="relative w-full h-screen overflow-hidden">
        {/* 地図（全画面） */}
        <Map
          defaultCenter={{ lat: 35.7100, lng: 139.7000 }}
          defaultZoom={12}
          disableDefaultUI
          mapId="history-map"
          style={{ width: "100%", height: "100%" }}
          onClick={() => setSelected(null)}
        >
          {/* 訪問済み場所のマーカー */}
          {places.map((place) => (
            <AdvancedMarker
              key={place.trip_id}
              position={{ lat: place.lat, lng: place.lng }}
              onClick={() => setSelected(place)}
            >
              <div className="w-9 h-9 bg-white border-2 border-blue-500 rounded-full flex items-center justify-center shadow-md hover:scale-110 transition-transform cursor-pointer">
                <span className="text-base">📍</span>
              </div>
            </AdvancedMarker>
          ))}

          {/* 選択中の場所の情報ウィンドウ */}
          {selected && (
            <InfoWindow
              position={{ lat: selected.lat, lng: selected.lng }}
              onCloseClick={() => setSelected(null)}
            >
              <div className="p-1 min-w-[160px]">
                <p className="font-bold text-gray-800 text-sm">
                  {selected.place_name}
                </p>
                <p className="text-xs text-gray-400 mt-0.5">
                  {new Date(selected.visited_at).toLocaleDateString("ja-JP")}
                </p>
                {selected.memory && (
                  <p className="text-xs text-gray-600 mt-2 bg-gray-50 rounded-lg p-2 leading-relaxed">
                    &quot;{selected.memory}&quot;
                  </p>
                )}
              </div>
            </InfoWindow>
          )}
        </Map>

        {/* 上部ヘッダー */}
        <div className="absolute top-4 left-4 right-4 flex items-center gap-3">
          <button
            onClick={() => router.push("/")}
            className="bg-white shadow-md rounded-xl px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
          >
            ← 戻る
          </button>
          <div className="bg-white shadow-md rounded-xl px-4 py-2">
            <p className="font-bold text-gray-800 text-sm">訪問した場所</p>
            <p className="text-xs text-gray-400">{places.length}件</p>
          </div>
        </div>
      </div>
    </APIProvider>
  );
}
