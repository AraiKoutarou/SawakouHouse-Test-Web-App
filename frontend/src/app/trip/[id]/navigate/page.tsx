"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  APIProvider,
  Map,
  AdvancedMarker,
  useMap,
} from "@vis.gl/react-google-maps";
import axios from "axios";

const API_BASE = "http://localhost:8080";

type Step = {
  instruction: string;
  distance_m: number;
  heading_deg: number;
};

type Position = {
  lat: number;
  lng: number;
};

// ── モックデータ（バックエンド連携前の開発用） ────────────────────
const MOCK_STEPS: Step[] = [
  { instruction: "北に向かって歩いてください", distance_m: 200, heading_deg: 0 },
  { instruction: "交差点を右折してください", distance_m: 150, heading_deg: 90 },
  { instruction: "左折してください", distance_m: 80, heading_deg: 270 },
  { instruction: "坂を上って正面の建物が目的地です", distance_m: 50, heading_deg: 0 },
];
// ─────────────────────────────────────────────────────────────────

// 地図の中心を現在地に追従させるコンポーネント
function MapContent({ userPos }: { userPos: Position }) {
  const map = useMap();

  useEffect(() => {
    if (map) map.panTo(userPos);
  }, [map, userPos]);

  return (
    <AdvancedMarker position={userPos}>
      {/* パルスエフェクト付き現在地マーカー */}
      <div className="relative flex items-center justify-center">
        <div className="absolute w-12 h-12 bg-blue-400 rounded-full animate-ping opacity-30" />
        <div className="absolute w-8 h-8 bg-blue-300 rounded-full opacity-50" />
        <div className="w-5 h-5 bg-blue-600 rounded-full border-2 border-white shadow-xl z-10" />
      </div>
    </AdvancedMarker>
  );
}

export default function NavigatePage() {
  const { id: tripId } = useParams<{ id: string }>();
  const router = useRouter();

  const [step, setStep] = useState<Step>(MOCK_STEPS[0]);
  const [stepIndex, setStepIndex] = useState(0);
  const [userPos, setUserPos] = useState<Position>({ lat: 35.6812, lng: 139.7671 });
  const [loading, setLoading] = useState(false);

  // 現在地を継続取得
  useEffect(() => {
    if (!navigator.geolocation) return;
    const watchId = navigator.geolocation.watchPosition(
      ({ coords }) =>
        setUserPos({ lat: coords.latitude, lng: coords.longitude }),
      (err) => console.error("Geolocation error:", err),
      { enableHighAccuracy: true }
    );
    return () => navigator.geolocation.clearWatch(watchId);
  }, []);

  const handleNext = async () => {
    setLoading(true);
    try {
      // TODO: バックエンド完成後に以下のコメントを外して切り替える
      // const res = await axios.post(`${API_BASE}/trips/${tripId}/next`, userPos);
      // if (res.data.arrived) {
      //   router.push(`/trip/${tripId}/arrived`);
      //   return;
      // }
      // setStep(res.data.step);

      // モック処理
      const next = stepIndex + 1;
      if (next >= MOCK_STEPS.length) {
        router.push(`/trip/${tripId}/arrived`);
      } else {
        setStepIndex(next);
        setStep(MOCK_STEPS[next]);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <APIProvider apiKey={process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY!}>
      <div className="relative w-full h-screen overflow-hidden">
        {/* 地図（全画面） */}
        <Map
          defaultCenter={userPos}
          defaultZoom={17}
          disableDefaultUI
          mapId="navigate-map"
          style={{ width: "100%", height: "100%" }}
        >
          <MapContent userPos={userPos} />
        </Map>

        {/* 上部：目的地非公開バッジ */}
        <div className="absolute top-4 left-0 right-0 flex justify-center pointer-events-none">
          <div className="bg-black/70 backdrop-blur-sm rounded-2xl px-5 py-2">
            <p className="text-white text-sm font-semibold tracking-widest">
              目的地：？？？
            </p>
          </div>
        </div>

        {/* 下部：ナビ指示カード */}
        <div className="absolute bottom-0 left-0 right-0 p-4">
          <div className="bg-white rounded-3xl shadow-2xl p-6">
            <div className="flex items-center gap-4 mb-5">
              {/* 方向矢印（heading_deg で回転） */}
              <div
                className="w-16 h-16 bg-blue-500 rounded-2xl flex items-center justify-center flex-shrink-0 shadow-md transition-transform duration-500"
                style={{ transform: `rotate(${step.heading_deg}deg)` }}
              >
                <svg
                  width="28"
                  height="28"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path d="M12 3L19 18L12 14L5 18L12 3Z" fill="white" />
                </svg>
              </div>

              <div className="flex-1 min-w-0">
                <p className="text-xl font-bold text-gray-800 leading-snug">
                  {step.instruction}
                </p>
                <p className="text-base text-blue-500 font-semibold mt-1">
                  {step.distance_m}m 先
                </p>
              </div>
            </div>

            <button
              onClick={handleNext}
              disabled={loading}
              className="w-full bg-gray-900 hover:bg-gray-700 disabled:bg-gray-300 text-white font-bold py-4 rounded-2xl transition-colors text-base"
            >
              {loading ? "読み込み中..." : "通過した / 到着した"}
            </button>
          </div>
        </div>
      </div>
    </APIProvider>
  );
}
