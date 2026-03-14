// lib/categories.ts: Google Places のタイプをアプリ独自のカテゴリーに変換。
export type CategoryInfo = {
  id: string;
  label: string;
  emoji: string;
  color: string;
  bg: string;
  border: string;
};

export const CATEGORIES: Record<string, CategoryInfo> = {
  gourmet: { id: 'gourmet', label: 'グルメ', emoji: '🍴', color: '#f97316', bg: '#fff7ed', border: '#fed7aa' },
  sightseeing: { id: 'sightseeing', label: '観光', emoji: '⛩️', color: '#ef4444', bg: '#fef2f2', border: '#fecaca' },
  nature: { id: 'nature', label: '自然', emoji: '🌲', color: '#22c55e', bg: '#f0fdf4', border: '#bbf7d0' },
  shopping: { id: 'shopping', label: '買物', emoji: '🛍️', color: '#ec4899', bg: '#fdf2f8', border: '#fbcfe8' },
  other: { id: 'other', label: 'その他', emoji: '📍', color: '#3b82f6', bg: '#eff6ff', border: '#bfdbfe' },
};

export function getCategoryFromType(types?: string[]): CategoryInfo {
  if (!types) return CATEGORIES.other;

  if (types.includes('restaurant') || types.includes('cafe') || types.includes('bar') || types.includes('food')) {
    return CATEGORIES.gourmet;
  }
  if (types.includes('tourist_attraction') || types.includes('museum') || types.includes('shrine') || types.includes('temple') || types.includes('place_of_worship')) {
    return CATEGORIES.sightseeing;
  }
  if (types.includes('park') || types.includes('natural_feature') || types.includes('campground')) {
    return CATEGORIES.nature;
  }
  if (types.includes('shopping_mall') || types.includes('store') || types.includes('department_store')) {
    return CATEGORIES.shopping;
  }
  
  return CATEGORIES.other;
}
