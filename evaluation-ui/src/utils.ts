import { format } from 'date-fns';
import type {Recipe} from './services/models/recipe.ts';

export function FormatBestDate(source:Recipe) {
    if(source.publishedDate) {
        const d = new Date(Date.parse(source.publishedDate));
        return format(d, "eee do MMM yyyy");
    } else {
        return "(no publish date)"
    }
}

export function FormatDate(dateString?: string) {
    if(dateString) {
        const d = new Date(Date.parse(dateString));
        return format(d, "eee do MMM yyyy");
    } else {
        return undefined;
    }
}