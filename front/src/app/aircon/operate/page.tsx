"use client";
import { useState, useEffect } from "react";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import Link from "next/link";
import Switch from "@mui/material/Switch";
import Slider from "@mui/material/Slider";

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

export default function Operate() {
  const [status, setStatus] = useState<AirconStatus>({
    power: "0",
    mode: "0",
    temp: "0",
    fan: "0",
    swing: "0",
    eco: "0",
    onTimer: "0",
    offTimer: "0",
  });
  const [error, setError] = useState<string | null>(null);

  const getStatus = async () => {
    try {
      const res = await fetch("/aircon/status-json", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const json = await res.json();
      if (res.status === 200) {
        console.log(json);
        setStatus(json);
      } else {
        setError(json.message);
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

  const operate = async (status: AirconStatus) => {
    try {
      console.log(JSON.stringify(status));

      const res = await fetch("/aircon/operate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(status),
      });
      const json = await res.json();
      if (res.status === 200) {
        console.log(json);
        getStatus();
      } else {
        setError(json.message);
        console.error("失敗しました");
      }
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <main className=" h-screen flex flex-col items-center justify-center">
      {error && (
        <div className=" fixed top-2 right-2 text-2xl rounded-xl p-2 px-6 bg-black shadow-lg  text-white">
          <p className="text-red-500">{error}</p>
        </div>
      )}
      <div className=" p-8 bg-white rounded-3xl shadow-lg w-2/3 h-2/3 flex flex-col items-center justify-start gap-8 relative">
        <h1 className="text-4xl">エアコンの操作</h1>
        <div className="flex flex-col items-center gap-2 w-2/3 justify-center">
          <div className="flex items-center gap-2 w-2/3 justify-center">
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">電源</h2>
              <Switch
                checked={status?.power === "1"}
                onChange={() => {
                  operate({
                    ...status,
                    power: status.power === "1" ? "0" : "1",
                  });
                }}
                inputProps={{ "aria-label": "controlled" }}
              />
            </div>
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">モード</h2>
              <select
                value={status.mode}
                className="p-2 border border-gray-300 rounded-lg w-40 cursor-pointer"
                onChange={(e) => {
                  operate({ ...status, mode: e.target.value });
                }}
              >
                <option value="0">暖房</option>
                <option value="1">除湿</option>
                <option value="2">冷房</option>
                <option value="3">送風</option>
              </select>
            </div>
            <div className="flex flex-col items-center gap-2 w-1/2">
              <h2 className="text-2xl">温度</h2>
              <Slider
                value={parseInt(status.temp)}
                onChange={(e, value) => {
                  operate({ ...status, temp: value.toString() });
                }}
                aria-label="Temperature"
                defaultValue={parseInt(status.temp)}
                valueLabelDisplay="auto"
                step={1}
                marks
                min={18}
                max={30}
              />
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
