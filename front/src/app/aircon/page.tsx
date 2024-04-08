"use client";
import { useState } from "react";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import Link from "next/link";

export default function ProxyLogin() {
  const [loggined, setLoggined] = useState(false);
  const [username, setUsername] = useState("");
  const checkLogin = async () => {
    try {
      const res = await fetch("/check-login", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (res.status === 200) {
        const json = await res.json();
        setLoggined(json.loggined);
        setUsername(json.username);
      } else {
        console.error("失敗しました");
      }
    } catch (e) {
      console.error(e);
    }
  };

  checkLogin();

  return (
    <main className=" h-screen flex flex-col items-center justify-center">
      <div className=" p-6 bg-white rounded-3xl shadow-lg w-2/3 h-2/3 flex flex-col items-center justify-center gap-4 relative">
        <LoginIcon className=" absolute top-2 right-2" loggined={loggined} />
        <div className="flex flex-col gap-2">
          <h1 className="text-4xl">管理画面</h1>
          <p>ようこそ、{loggined ? username : "ゲスト"} さん</p>
        </div>

        <div className="flex items-center gap-4 w-2/3 justify-center border-4 p-4 rounded-2xl shadow-md h-3/4">
          <a
            className="p-2 rounded-2xl shadow-lg px-8 hover:bg-blue-500 bg-blue-400 w-2/3 aspect-square text-center text-white text-3xl"
            href="/aircon/status"
          >
            エアコンの状態をみる
          </a>
          {loggined && username === "admin" ? (
            <a
              className="p-2 rounded-2xl shadow-lg px-8 hover:bg-blue-500 bg-blue-400 w-2/3 aspect-square text-center text-white text-3xl"
              href="/aircon/operate"
            >
              エアコンを操作する
            </a>
          ) : (
            <a
              className="p-2 rounded-2xl shadow-lg px-8  bg-gray-400 w-2/3 aspect-square text-center text-white text-3xl cursor-not-allowed"
              href=""
            >
              エアコンを操作する
            </a>
          )}
        </div>
      </div>
    </main>
  );
}

function LoginIcon({
  className,
  loggined,
}: {
  className?: string;
  loggined: boolean;
}) {
  const [showUserMenu, setShowUserMenu] = useState(false);
  console.log(loggined);
  return (
    <div className={className}>
      <div
        className="flex flex-col items-end justify-center text-slate-600 rounded-full gap-4"
        onClick={() => setShowUserMenu(!showUserMenu)}
      >
        <AccountCircleIcon
          sx={{ fontSize: 80 }}
          className="shadow-md rounded-full cursor-pointer"
        />
        <ul
          className="text-xl w-48 bg-slate-100 rounded-xl shadow-md flex flex-col items-center justify-center gap-6 p-2 text-slate-600"
          style={{ display: showUserMenu ? "block" : "none" }}
        >
          {!loggined ? (
            <li className="cursor-pointer hover:text-slate-800 hover:bg-slate-200 rounded-full m-2">
              <a href="/login" className="w-100 block">
                ログイン
              </a>
            </li>
          ) : (
            <li className="cursor-pointer hover:text-slate-800 hover:bg-slate-200 rounded-full m-2">
              <a href="/logout" className="w-100 block">
                ログアウト
              </a>
            </li>
          )}
          <li className="cursor-pointer hover:text-slate-800 hover:bg-slate-200 rounded-full m-2">
            <a href="/user" className="w-100 block">
              登録情報変更
            </a>
          </li>
        </ul>
      </div>
    </div>
  );
}
