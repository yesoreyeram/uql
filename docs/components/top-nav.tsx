import Link from "next/link";
import { Logo } from "@/components/logo";
import { NavigationMenu, NavigationMenuItem, NavigationMenuList } from "@/components/ui/navigation-menu";

export const menuItems = [
  {
    name: "docs",
    href: "/docs",
  },
  {
    name: "playground",
    href: "/playground",
  },
  {
    name: "examples",
    href: "/examples",
  },
  {
    name: "github",
    href: "https://github.com/yesoreyeram/uql",
    position: "right",
  },
];

export const TopNav = () => {
  return (
    <div className="flex justify-center bg-white border-b-1 border-black sticky top-0">
      <div className="w-6xl">
        <div className="flex flex-row justify-between">
          <NavigationMenu>
            <NavigationMenuList className="gap-0">
              <NavigationMenuItem className="py-1 px-2">
                <Logo />
              </NavigationMenuItem>
              {menuItems
                .filter((item) => item.position !== "right")
                .map((item) => (
                  <Link href={item.href} key={item.name}>
                    <MenuItem name={item.name} />
                  </Link>
                ))}
            </NavigationMenuList>
          </NavigationMenu>
          <NavigationMenu>
            <NavigationMenuList className="gap-5">
              {menuItems
                .filter((item) => item.position === "right")
                .map((item) => (
                  <Link href={item.href} key={item.name} target="_blank">
                    <MenuItem name={item.name} />
                  </Link>
                ))}
            </NavigationMenuList>
          </NavigationMenu>
        </div>
      </div>
    </div>
  );
};

const MenuItem = (item: { name: string }) => {
  return (
    <NavigationMenuItem key={item.name} className="hover:text-white hover:bg-black text-center py-2 px-4">
      {item.name}
    </NavigationMenuItem>
  );
};
