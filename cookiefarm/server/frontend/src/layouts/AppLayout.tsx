import { useEffect, useState } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Sidebar } from "@cloudflare/kumo/components/sidebar";
import { Switch } from "@cloudflare/kumo/components/switch";
import { Text } from "@cloudflare/kumo/components/text";
import { z } from "zod";
import {
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
import { useTheme } from "@/hooks/useTheme";

const apiStatusSchema = z.object({
  message: z.string(),
  time: z.string(),
});

const configHeaderSchema = z.object({
  configured: z.boolean().default(false),
  server: z.object({
    protocol: z.string().trim().default(""),
  }),
});

const navigationItems = [
  { href: "/", label: "Dashboard", icon: ChartBar },
  { href: "/flags", label: "Flags", icon: Flag },
  { href: "/config", label: "Config", icon: Gear },
] as const;

export function AppLayout() {
  const auth = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const { mode, toggleMode } = useTheme();
  const [competitionName, setCompetitionName] = useState("CookieFarm");
  const [status, setStatus] = useState<"connecting" | "open" | "closed" | "error">(
    "connecting",
  );

  useEffect(() => {
    void apiFetch("/config", {}, configHeaderSchema)
      .then((config) => {
        const suffix = config.server.protocol ? ` / ${config.server.protocol}` : "";
        setCompetitionName(`CookieFarm${suffix}`);
      })
      .catch(() => {
        setCompetitionName("CookieFarm");
      });
  }, []);

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
    <Sidebar.Provider defaultOpen variant="floating" collapsible="icon">
      {/* ① Use Tailwind grid directly — no custom class needed */}
      <div className="grid min-h-screen grid-cols-[auto_minmax(0,1fr)]">
        <Sidebar className="border-r border-kumo-line/60 bg-kumo-base/95">
          <Sidebar.Header className="px-3 py-4">
            <div className="flex items-center gap-3">
              <div className="grid size-9 shrink-0 place-items-center rounded-xl bg-kumo-brand text-white">
                <Flag size={18} weight="fill" />
              </div>
              {/* ② Show dynamic name here too, not just the top bar */}
              <div className="min-w-0">
                <div className="truncate font-semibold text-kumo-fg-primary">
                  {competitionName}
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
            <span className="truncate text-sm text-kumo-fg-secondary">Dense mode</span>
          </Sidebar.Footer>

          <Sidebar.Rail />
        </Sidebar>

        {/* ④ Right panel: flex column so header is sticky within this column */}
        <div className="flex min-w-0 flex-col">
          <header className="sticky top-0 z-30 border-b border-kumo-line/70 bg-kumo-overlay/90 px-4 py-3 backdrop-blur md:px-6">
            <div className="mx-auto flex max-w-[1600px] items-center justify-between gap-4">
              <div className="min-w-0">
                <Text size="sm" variant="secondary">
                  Console
                </Text>
                <div className="truncate font-semibold text-kumo-fg-primary">
                  {competitionName}
                </div>
              </div>

              <div className="flex items-center gap-4">
                <LiveDot status={status} showLabel />
                <div className="flex items-center gap-2">
                  {mode === "dark" ? (
                    <Moon size={16} className="text-kumo-fg-secondary" />
                  ) : (
                    <Sun size={16} className="text-kumo-fg-secondary" />
                  )}
                  <Switch
                    label={<span className="sr-only">Toggle dark mode</span>}
                    checked={mode === "dark"}
                    onCheckedChange={toggleMode}
                    aria-label="Toggle dark mode"
                  />
                </div>
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

          <main className="mx-auto min-w-0 w-full max-w-[1600px] p-4 md:p-6">
            <Outlet />
          </main>
        </div>
      </div>
    </Sidebar.Provider>
  );
}
