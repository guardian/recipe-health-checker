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