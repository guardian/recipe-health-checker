import {Recipe} from "./models/recipe.ts";

const ApiBaseUrl = "https://recipes.guardianapis.com";
const InternalCache:Map<string, Recipe> = new Map();

export async function addToCache(recipeId:string, overwrite:boolean=false) {
    if(!overwrite && InternalCache.has(recipeId)) return InternalCache.get(recipeId);

    const url = `${ApiBaseUrl}/api/content/by-uid/${recipeId}`;
    const rawContent = await fetch(url);
    if(rawContent.status===200) {
        const rawJson = await rawContent.json();
        const parsed = Recipe.parse(rawJson);
        InternalCache.set(recipeId, parsed);
        return parsed;
    } else {
        const responseData = await rawContent.text();
        console.error(`Unable to retrieve recipe with id ${recipeId}: ${rawContent.status}`,  responseData);
        throw new Error(`Could not retrieve recipe: ${rawContent.status}`);
    }
}

export function removeFromCache(recipeId:string) {
    InternalCache.delete(recipeId);
}

export async function retrieveFromCache(recipeId:string) {
    const cachedCopy = InternalCache.get(recipeId);
    if(cachedCopy) {
        return cachedCopy;
    } else {
        return addToCache(recipeId);
    }
}

export async function cacheEverythingInList(idList:Array<string>) {
    for(const recipeId of idList) {
        try {
            await addToCache(recipeId);
        } catch(err) {
            console.warn(`Could not load ID ${recipeId}: ${err}`);
        }
    }
}

export function retrieveCacheAsRecord() {
    return Object.fromEntries(InternalCache);
}