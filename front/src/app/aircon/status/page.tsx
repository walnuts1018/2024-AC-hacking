"use client";
import { useState, useEffect } from "react";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import Link from "next/link";

type AirconStatus = {
  power: string;
  mode: string;
  temp: string;
  fan: string;
  swing: string;
  eco: string;
  onTimer: string;
  offTimer: string;
};

const fetchInterval = 3 * 1000;

export default function AirconStatus() {
  const [status, setStatus] = useState<AirconStatus | null>(null);
  const getStatus = async () => {
    try {
      const res = await fetch("/aircon/status-json", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (res.status === 200) {
        const json = await res.json();
        setStatus(json);
      } else {
        console.error("失敗しました");
      }
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    getStatus();
    setInterval(async () => {
      getStatus();
    }, fetchInterval);
    return () => {
      clearInterval(fetchInterval);
    };
  }, []);

  return (
    <main className=" h-screen flex flex-col items-center justify-center">
      <div className=" p-8 bg-white rounded-3xl shadow-lg w-2/3 h-2/3 flex flex-col items-center justify-start gap-8 relative">
        <h1 className="text-4xl">エアコンの状態</h1>
        <p className="text-xl">
          これは過去の操作履歴から取得したものであり、実際の状態と不整合がある場合があります。
        </p>
        <div className="flex flex-col items-center gap-2 w-2/3 justify-center">
          <div className="flex items-center gap-2 w-2/3 justify-center">
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">電源</h2>
              <p>{status?.power === "1" ? "ON" : "OFF"}</p>
            </div>
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">モード</h2>
              <p>
                {status?.mode === "0"
                  ? "暖房"
                  : status?.mode === "1"
                  ? "除湿"
                  : status?.mode === "2"
                  ? "冷房"
                  : status?.mode === "3"
                  ? "送風"
                  : "不明"}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2 w-2/3 justify-center">
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">温度</h2>
              <p>{status?.temp}度</p>
            </div>
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">風量</h2>
              <p>
                {status?.fan === "0"
                  ? "自動"
                  : status?.fan === "1"
                  ? "弱"
                  : status?.fan === "2"
                  ? "中"
                  : "強"}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2 w-2/3 justify-center">
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">風向</h2>
              <p>{status?.swing === "0" ? "固定" : "自動"}</p>
            </div>
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">エコモード</h2>
              <p>{status?.eco === "1" ? "ON" : "OFF"}</p>
            </div>
          </div>
          <div className="flex items-center gap-2 w-2/3 justify-center">
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">ONタイマー</h2>
              <p>{status?.onTimer}時間後</p>
            </div>
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">OFFタイマー</h2>
              <p>{status?.offTimer}時間後</p>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
