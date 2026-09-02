import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { BrowserRouter } from "react-router-dom";

import App from "@/app/App";
import { store } from "@/app/store";
import { ThemeProvider } from "@/shared/theme/ThemeProvider";
import InstallPrompt from "@/shared/pwa/InstallPrompt";
import { registerServiceWorker } from "@/shared/pwa/registerSW";
import ErrorBoundary from "@/shared/ui/ErrorBoundary";
import { ToastProvider } from "@/shared/ui/ToastProvider";

import "./index.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ErrorBoundary>
      <Provider store={store}>
        <BrowserRouter>
          <ThemeProvider>
            <ToastProvider>
              <App />
              <InstallPrompt />
            </ToastProvider>
          </ThemeProvider>
        </BrowserRouter>
      </Provider>
    </ErrorBoundary>
  </StrictMode>,
);

registerServiceWorker();
