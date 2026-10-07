// src/components/layouts/types.ts
export interface MenuItem {
    path: string
    label: string
    icon: string
}

export interface MenuGroup {
    id: string
    label: string
    icon: string
    items: MenuItem[]
}