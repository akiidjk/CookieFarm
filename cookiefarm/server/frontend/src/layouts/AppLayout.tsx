import { useState } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Sidebar } from "@cloudflare/kumo/components/sidebar";
import { Text } from "@cloudflare/kumo/components/text";
import { z } from "zod";
import {
  Bug,
  ChartBar,
  Flag,
  Gear,
  SignOut,
  Moon,
  SidebarSimple,
  Sun,
} from "@phosphor-icons/react";
import { Outlet, useLocation, useNavigate } from "react-router";
import { apiFetch } from "@/api/client";
import { useAuth } from "@/features/auth/AuthProvider";
import { LiveDot } from "@/components/LiveDot";
import { useInterval } from "@/hooks/useInterval";
import { CookieIcon, HouseIcon, HouseSimpleIcon } from "@phosphor-icons/react/dist/ssr";

const apiStatusSchema = z.object({
  message: z.string(),
  time: z.string(),
});


const navigationItems = [
  { href: "/", label: "Dashboard", icon: HouseIcon },
  { href: "/charts", label: "Charts", icon: ChartBar },
  { href: "/flags", label: "Flags", icon: Flag },
  { href: "/exploits", label: "Exploits", icon: Bug },
  { href: "/config", label: "Config", icon: Gear },
] as const;

export function AppLayout() {
  const auth = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [status, setStatus] = useState<"connecting" | "open" | "closed" | "error">(
    "connecting",
  );

  useInterval(
    () => {
      void apiFetch("/", {}, apiStatusSchema)
        .then((response) => {
          setStatus(response.message ? "open" : "error");
        })
        .catch(() => {
          setStatus("closed");
        });
    },
    15000,
    { immediate: true },
  );

  return (
    <Sidebar.Provider defaultOpen variant="inset" collapsible="icon">
      {/* ① Use Tailwind grid directly — no custom class needed */}
      <div className="absolute inset-0 flex overflow-hidden bg-kumo-canvas">
        <Sidebar className="border-r border-kumo-line/60 bg-kumo-base/95">
          <Sidebar.Header className="px-3 py-4">
            <div className="flex items-center gap-3">
              <img src="/images/logo.png" alt="CookieFarm Logo" width={24} height={24} className="block" />
              {/* ② Show dynamic name here too, not just the top bar */}
              <div className="min-w-0">
                <div className="truncate font-semibold text-kumo-fg-primary">
                  Cookiefarm
                </div>
                <Text size="sm" variant="secondary">
                  Operator Console
                </Text>
              </div>
            </div>
          </Sidebar.Header>

          <Sidebar.Content className="flex-1 overflow-y-auto px-2 pb-3">
            <Sidebar.Group>
              <Sidebar.GroupLabel>Operations</Sidebar.GroupLabel>
              <Sidebar.Menu>
                {navigationItems.map((item) => (
                  <Sidebar.MenuButton
                    key={item.href}
                    icon={item.icon}
                    active={
                      item.href === "/"
                        ? location.pathname === item.href
                        : location.pathname.startsWith(item.href)
                    }
                    onClick={() => navigate(item.href)}
                  >
                    {item.label}
                  </Sidebar.MenuButton>
                ))}
              </Sidebar.Menu>
            </Sidebar.Group>
          </Sidebar.Content>

          {/* ③ Footer: trigger on left, label only visible when expanded */}
          <Sidebar.Footer className="flex items-center gap-2 px-3 py-4">
            <Sidebar.Trigger aria-label="Toggle sidebar">
              <SidebarSimple size={18} />
            </Sidebar.Trigger>
          </Sidebar.Footer>

          <Sidebar.Rail />
        </Sidebar>

        {/* ④ Right panel: flex column so header is sticky within this column */}
        <div className="flex flex-1 min-w-0 flex-col overflow-y-auto">
          <header className="sticky top-0 z-30 border-b border-kumo-line/70 bg-kumo-overlay/90 px-4 py-3 backdrop-blur md:px-6">
            <div className="mx-auto flex items-center justify-between gap-4">
              <div className="min-w-0">
                <Text size="sm" variant="secondary">
                  Console
                </Text>
                <div className="truncate font-semibold text-kumo-fg-primary">
                  CookieFarm
                </div>
              </div>

              <div className="flex items-center gap-4">
                <LiveDot status={status} showLabel />
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => {
                    void auth.logout().then(() => {
                      navigate("/login", { replace: true });
                    });
                  }}
                >
                  <SignOut size={16} />
                  Sign out
                </Button>
              </div>
            </div>
          </header>

          <main className="mx-auto min-w-0 w-full p-4 md:p-6">
            <Outlet />
          </main>
        </div>
      </div>
    </Sidebar.Provider>
  );
}
