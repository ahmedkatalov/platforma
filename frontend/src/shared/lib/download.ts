import { tokenStorage } from "@/shared/api/tokenStorage";

// Скачивает файл с бэкенда: обычная ссылка не подойдёт — нужен заголовок авторизации.
export async function downloadFile(path: string, fallbackName: string): Promise<void> {
  const response = await fetch(path, {
    headers: { Authorization: `Bearer ${tokenStorage.access() ?? ""}` },
  });

  if (!response.ok) {
    throw new Error("Не удалось скачать файл");
  }

  // Имя файла приходит в Content-Disposition.
  const disposition = response.headers.get("Content-Disposition") ?? "";
  const match = disposition.match(/filename="?([^"]+)"?/);
  const filename = match?.[1] ?? fallbackName;

  const blob = await response.blob();
  const url = URL.createObjectURL(blob);

  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();

  URL.revokeObjectURL(url);
}
