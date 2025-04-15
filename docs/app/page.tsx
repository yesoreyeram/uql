import Link from "next/link";
import { Logo } from "@/components/logo";
import { menuItems } from "@/components/top-nav";

export default function Home() {
  return (
    <div className="grid grid-rows-[20px_1fr_20px] min-h-screen items-center justify-items-center p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)] bg-white text-black">
      <main className="flex flex-col row-start-2 items-center sm:items-start justify-content">
        <div className="flex flex-col gap-2 border-1 py-4 px-6 border-black bg-white w-2xl">
          <Logo />
          <h2>Unified Query Language</h2>
        </div>
        <ul className="flex gap-0 w-full border-black border-b-1 border-x-1 bg-white text-black">
          {menuItems.map((item) => (
            <Link href={item.href} key={item.name} className="hover:bg-black hover:text-white w-1/3 text-center py-1">
              <li>{item.name}</li>
            </Link>
          ))}
        </ul>
      </main>
    </div>
  );
}
