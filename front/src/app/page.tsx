import Image from "next/image";

export default function Home() {
  return (
    <main className=" h-screen flex flex-col items-center justify-center">
      <h1>エアコン操作サイト</h1>
      <p>エアコンを操作するサイトです</p>
      <Image src="/aircon.jpg" alt="エアコン" width={500} height={500} />
    </main>
  );
}
