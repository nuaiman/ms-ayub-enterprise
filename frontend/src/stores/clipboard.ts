import { ref } from "vue";
import { defineStore } from "pinia";
import { push } from "notivue";

export const useClipboardStore = defineStore('clipboard', () => {
    const copied = ref(false);
    let timeoutId: ReturnType<typeof setTimeout> | null = null;

    const copyToClipboard = async (text: string): Promise<boolean> => {
        try {
            await navigator.clipboard.writeText(text);
            showCopiedFeedback();
            return true;
        } catch {
            // Fallback
            try {
                const textarea = document.createElement('textarea');
                textarea.value = text;
                textarea.style.position = 'fixed';
                textarea.style.opacity = '0';
                document.body.appendChild(textarea);
                textarea.select();
                document.execCommand('copy');
                document.body.removeChild(textarea);
                showCopiedFeedback();
                return true;
            } catch {
                push.error('Failed to copy to clipboard');
                return false;
            }
        }
    };

    const showCopiedFeedback = () => {
        copied.value = true;
        push.success('Copied to clipboard!');
        if (timeoutId) clearTimeout(timeoutId);
        timeoutId = setTimeout(() => {
            copied.value = false;
        }, 2000);
    };

    const resetCopied = () => {
        copied.value = false;
        if (timeoutId) {
            clearTimeout(timeoutId);
            timeoutId = null;
        }
    };

    return {
        copied,
        copyToClipboard,
        resetCopied,
    };
});