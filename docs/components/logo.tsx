import Link from "next/link";

export const Logo = () => {
  return (
    <Link href="/welcome">
      <h1 className="gap-0 font-bold">
        <span className="mr-1">{`{U}`}</span>
        <span className="">UQL</span>
      </h1>
    </Link>
  );
};
