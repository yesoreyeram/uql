import { TopNav } from "@/components/top-nav";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "UQL",
  description: "UQL",
};

type Props = Readonly<{ children: React.ReactNode }>;

export default function Layout({ children }: Props) {
  return (
    <div>
      <TopNav />
      <div className="px-6 py-2">
        <div className="flex justify-center">
          <div className="w-6xl">
            <div className="mdx">{children}</div>
          </div>
        </div>
      </div>
    </div>
  );
}
