// Загрузка движка v86 (эмулятор x86 на WebAssembly) — реальный Linux в браузере.
// Библиотека и образ большие (~10 МБ), поэтому грузятся только по требованию,
// когда студент открывает песочницу и нажимает «Запустить».

let loading: Promise<V86Ctor> | null = null;

export type V86Ctor = new (options: Record<string, unknown>) => V86Instance;

export type V86Instance = {
  add_listener(event: string, handler: (data: unknown) => void): void;
  remove_listener?(event: string, handler: (data: unknown) => void): void;
  serial0_send(data: string): void;
  destroy?(): Promise<void> | void;
  save_state(): Promise<ArrayBuffer>;
  restore_state(state: ArrayBuffer): Promise<void>;
};

// Пути к самостоятельно захостированным ассетам (frontend/public/v86).
export const V86_BASE = "/v86";
export const V86_ASSETS = {
  wasm_path: `${V86_BASE}/v86.wasm`,
  bios: `${V86_BASE}/seabios.bin`,
  vga_bios: `${V86_BASE}/vgabios.bin`,
  cdrom: `${V86_BASE}/linux4.iso`,
};

export function loadV86(): Promise<V86Ctor> {
  const w = window as unknown as { V86?: V86Ctor };
  if (w.V86) return Promise.resolve(w.V86);
  if (loading) return loading;

  loading = new Promise<V86Ctor>((resolve, reject) => {
    const script = document.createElement("script");
    script.src = `${V86_BASE}/libv86.js`;
    script.async = true;
    script.onload = () => {
      const ctor = (window as unknown as { V86?: V86Ctor }).V86;
      if (ctor) resolve(ctor);
      else reject(new Error("Движок v86 загрузился, но конструктор не найден"));
    };
    script.onerror = () => {
      loading = null;
      reject(new Error("Не удалось загрузить движок терминала (v86)"));
    };
    document.head.appendChild(script);
  });
  return loading;
}
