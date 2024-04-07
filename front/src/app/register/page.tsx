"use client";
import { useState } from "react";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  const login = async () => {
    if (!username) {
      setError("ユーザ名を入力してください");
      return;
    }

    if (!password) {
      setError("パスワードを入力してください");
      return;
    }
    try {
      const res = await fetch("/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });
      if (res.status === 200) {
        window.location.href = "/aircon";
      } else {
        setError("登録に失敗しました");
      }
    } catch (e) {
      setError("登録に失敗しました");
    }
  };

  return (
    <main className=" h-screen flex flex-col items-center justify-center">
      <div className=" p-4 bg-white rounded-3xl shadow-lg w-2/3 h-2/3 flex flex-col items-center justify-center gap-16">
        <h1 className="text-4xl">ユーザー登録</h1>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            login();
          }}
          className="flex flex-col items-center gap-4"
          id="login-form"
        >
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="p-2 border border-gray-300 rounded-lg w-80"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="p-2 border border-gray-300 rounded-lg w-80"
          />
        </form>
        <button
          type="submit"
          form="login-form"
          className="p-2 px-4 bg-blue-500 text-white rounded-full"
        >
          登録
        </button>
        <p className="text-black p-2 rounded-lg">{error}</p>
      </div>
    </main>
  );
}
