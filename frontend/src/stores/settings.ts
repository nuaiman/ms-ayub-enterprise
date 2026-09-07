import { ref, computed, watch } from "vue";
import { defineStore } from "pinia";
import api from "@/utils/axios";
import type { Theme } from "@/types/settings";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";

export const useSettingsStore = defineStore('settings', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader();

    // STATE
    const theme = ref<Theme>((localStorage.getItem('theme') as Theme) || 'light');

    // COMPUTED
    const isDark = computed(() => theme.value === 'dark');
    const themeLabel = computed(() =>
        theme.value === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'
    );

    // ACTIONS
    const applyTheme = (value: Theme) => {
        theme.value = value;
        document.documentElement.setAttribute('data-theme', value);
        localStorage.setItem('theme', value);
    };

    const toggleTheme = () => {
        applyTheme(theme.value === 'dark' ? 'light' : 'dark');
    };

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

    // WATCH
    watch(theme, (newTheme) => {
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
    }, { immediate: true });

    return {
        theme,
        isDark,
        themeLabel,
        applyTheme,
        toggleTheme,
        downloadBackup,
    };
});