import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router";
import { Toasty } from "@cloudflare/kumo/components/toast";
import { startMocking } from "@/api/mock";
import { AuthProvider } from "@/features/auth/AuthProvider";
import { ThemeProvider } from "@/hooks/useTheme";
import { router } from "@/router";
import "@/styles/main.css";

async function bootstrap() {
  await startMocking();

  const root = document.getElementById("root");
  if (!root) {
    throw new Error("Root element not found");
  }

  createRoot(root).render(
    <StrictMode>
      <ThemeProvider>
        <AuthProvider>
          <Toasty>
            <RouterProvider router={router} />
          </Toasty>
        </AuthProvider>
      </ThemeProvider>
    </StrictMode>,
  );
}

void bootstrap();
