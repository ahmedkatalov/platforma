import { useRef, useState } from "react";

import {
  useDeleteAssetMutation,
  useGetAssetsQuery,
  useUploadAssetMutation,
} from "@/features/certificates/api/certificatesApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { Button, EmptyState, Modal, Spinner } from "@/shared/ui";
import { IconPlus, IconTrash } from "@/shared/ui/icons";
import { useToast } from "@/shared/ui/ToastProvider";

// Библиотека картинок для уроков: загрузка, выбор и удаление.
// По клику отдаёт готовую markdown-вставку.
export default function AssetPicker({
  open,
  onClose,
  onPick,
}: {
  open: boolean;
  onClose: () => void;
  onPick: (markdown: string) => void;
}) {
  const { data: assets = [], isLoading } = useGetAssetsQuery(undefined, { skip: !open });
  const [uploadAsset, { isLoading: uploading }] = useUploadAssetMutation();
  const [deleteAsset] = useDeleteAssetMutation();

  const inputRef = useRef<HTMLInputElement>(null);
  const [selected, setSelected] = useState<string | null>(null);
  const toast = useToast();

  const upload = async (file: File) => {
    try {
      const asset = await uploadAsset(file).unwrap();
      setSelected(asset.id);
      toast.success("Картинка загружена");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось загрузить файл"));
    }
  };

  const remove = async (id: string) => {
    if (!window.confirm("Удалить картинку? Она пропадёт во всех уроках, где использовалась.")) return;
    try {
      await deleteAsset(id).unwrap();
      if (selected === id) setSelected(null);
      toast.success("Картинка удалена");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const insert = () => {
    const asset = assets.find((item) => item.id === selected);
    if (!asset) return;
    onPick(`![${asset.original || "изображение"}](${asset.url})`);
    onClose();
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title="Картинки уроков"
      width="42rem"
      footer={
        <>
          <Button onClick={onClose}>Отмена</Button>
          <Button variant="primary" onClick={insert} disabled={!selected}>
            Вставить в текст
          </Button>
        </>
      }
    >
      <input
        ref={inputRef}
        type="file"
        accept="image/png,image/jpeg,image/gif,image/webp,image/svg+xml"
        className="hidden"
        onChange={(e) => {
          const file = e.target.files?.[0];
          if (file) void upload(file);
          e.target.value = "";
        }}
      />

      <div className="mb-4 flex items-center justify-between gap-3">
        <p className="text-sm text-muted">
          PNG, JPG, GIF, WebP или SVG — до 5 МБ. Ссылка вставится как markdown.
        </p>
        <Button
          variant="primary"
          icon={<IconPlus size={16} />}
          onClick={() => inputRef.current?.click()}
          loading={uploading}
        >
          Загрузить
        </Button>
      </div>

      {isLoading ? (
        <div className="grid place-items-center py-12 text-accent">
          <Spinner size={28} />
        </div>
      ) : assets.length === 0 ? (
        <EmptyState
          title="Картинок пока нет"
          description="Загрузите схему или скриншот — они появятся в этой библиотеке"
        />
      ) : (
        <ul className="grid grid-cols-2 gap-3 sm:grid-cols-3">
          {assets.map((asset) => (
            <li key={asset.id}>
              <button
                type="button"
                onClick={() => setSelected(asset.id)}
                className={`group relative block w-full overflow-hidden rounded-[var(--radius-md)] border transition-colors ${
                  selected === asset.id
                    ? "border-[var(--accent)] ring-2 ring-[var(--accent-soft)]"
                    : "border-line hover:border-[var(--accent-border)]"
                }`}
              >
                <span className="grid h-24 place-items-center bg-surface-2">
                  <img
                    src={asset.url}
                    alt={asset.original}
                    className="max-h-24 w-full object-contain"
                    loading="lazy"
                  />
                </span>
                <span className="block truncate px-2 py-1.5 text-left text-[11px] text-muted">
                  {asset.original || asset.filename}
                </span>
              </button>

              <button
                type="button"
                className="mt-1 flex w-full items-center justify-center gap-1 text-[11px] text-faint transition-colors hover:text-danger"
                onClick={() => remove(asset.id)}
              >
                <IconTrash size={12} />
                удалить
              </button>
            </li>
          ))}
        </ul>
      )}
    </Modal>
  );
}
