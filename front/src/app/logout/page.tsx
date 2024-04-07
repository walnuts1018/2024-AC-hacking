"use client";
import { useState } from "react";

export default function Logout() {
  const logout = async () => {
    try {
      const res = await fetch("/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (res.status === 200) {
      }
    } catch (e) {
      console.error(e);
    }
  };

  logout();

  return (
    <main className=" h-screen flex flex-col items-center justify-center">
      <div className=" p-4 bg-white rounded-3xl shadow-lg w-2/3 h-2/3 flex flex-col items-center justify-center gap-16">
        <h1 className="text-4xl">ログアウト</h1>
        <p className="text-black p-2 rounded-lg">ログアウトしました</p>
        <a
          className="p-2 rounded-full shadow-lg px-8 hover:bg-gray-100"
          href="/login"
        >
          再度ログイン
        </a>
      </div>
    </main>
  );
}
