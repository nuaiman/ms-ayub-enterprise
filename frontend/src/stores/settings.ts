// src/stores/settings.ts
import { defineStore } from "pinia";
import api from "@/utils/axios";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";

export const useSettingsStore = defineStore('settings', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader();

    const downloadBackup = async (): Promise<boolean> => {
        displayLoader();

        try {
            const res = await api.get('/backup', {
                responseType: 'blob',
            });

            let filename = 'backup.zip';
            const disposition = res.headers['content-disposition'];
            if (disposition) {
                const match = disposition.match(/filename="?([^"]+)"?/);
                if (match) filename = match[1];
            }

            const blob = new Blob([res.data]);
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = filename;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);

            push.success('Backup downloaded successfully!');
            return true;
        } catch (err) {
            console.error(err);
            push.error('Failed to download backup');
            return false;
        } finally {
            destroyLoader();
        }
    };

    return {
        downloadBackup,
    };
});