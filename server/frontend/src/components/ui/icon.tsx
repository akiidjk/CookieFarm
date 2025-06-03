import { Icon as IconifyIcon } from "@iconify/react";
import { cn } from "@/lib/utils";

type IconProps = {
  name: string;
  className?: string;
  size?: number | string;
  color?: string;
  onClick?: () => void;
};

export function Icon({
  name,
  className,
  size,
  color,
  onClick,
  ...props
}: IconProps & Omit<React.ComponentProps<typeof IconifyIcon>, "icon">) {
  return (
    <IconifyIcon
      icon={name}
      className={cn("flex-shrink-0", className)}
      width={size}
      height={size}
      color={color}
      onClick={onClick}
      {...props}
    />
  );
}

export function getIconName(category: string, name: string): string {
  return `${category}:${name}`;
}