export const CURRENCY_SYMBOL = '৳';

export const formatCurrency = (amount: number | null | undefined): string => {
    if (amount === null || amount === undefined || isNaN(amount)) return '—';
    return `${CURRENCY_SYMBOL} ${amount.toFixed(2)}`;
};

export const parseCurrency = (value: string): number | null => {
    const cleaned = value.replace(/[^0-9.]/g, '');
    const parsed = parseFloat(cleaned);
    return isNaN(parsed) ? null : parsed;
};

export const formatCurrencyCompact = (amount: number | null | undefined): string => {
    if (amount === null || amount === undefined || isNaN(amount)) return '—';
    if (amount >= 1000000) return `${CURRENCY_SYMBOL} ${(amount / 1000000).toFixed(1)}M`;
    if (amount >= 1000) return `${CURRENCY_SYMBOL} ${(amount / 1000).toFixed(1)}K`;
    return `${CURRENCY_SYMBOL} ${amount.toFixed(2)}`;
};