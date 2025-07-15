import { z } from "zod";

const Range = z.object({
    min: z.number().nullable(),
    max: z.number().nullable(),
});

const RecipeImage = z.object({
    url: z.string(),
    mediaId: z.string(),
    cropId: z.string(),
    source: z.string().optional(),
    photographer: z.string().optional(),
    imageType: z.string().optional(),
    caption: z.string().optional().nullable(),
    mediaApiUri: z.string().optional(),
    displayCredit: z.boolean().optional(),
    width: z.number().optional(),
    height: z.number().optional(),
});

const Ingredient = z.object({
    name: z.string(),
    ingredientId: z.string().optional(),
    amount: Range.optional(),
    unit: z.string().optional(),
    prefix: z.string().optional(),
    suffix: z.string().optional(),
    text: z.string().optional(),
    optional: z.boolean().default(false),
    emptyAmountIsOk: z.boolean().optional(),
});

const IngredientsGroup = z.object({
    recipeSection: z.string(),
    ingredientsList: z.array(Ingredient),
});

const Timing = z.object({
    qualifier: z.string(),
    durationInMins: Range,
    text: z.string().optional(),
});

const Instruction = z.object({
    description: z.string(),
    images: z.array(RecipeImage).optional(),
});

const Serves = z.object({
    amount: Range,
    unit: z.string().optional(),
    text: z.string().optional(),
});

const CommerceCTA = z.object({
    territory: z.string(),
    sponsorName: z.string(),
    url: z.string(),
});

// Accept string or object as contributor
const Contributor = z.union([
    z.string(),
    z.object({
        name: z.string(),
        role: z.string().optional(),
    }),
]);

export const Recipe = z.object({
    id: z.string(),
    composerId: z.string().optional(),
    canonicalArticle: z.string(),
    title: z.string(),
    description: z.string(),
    featuredImage: RecipeImage.optional(),
    previewImage: RecipeImage.optional(),
    contributors: z.array(Contributor).optional(), // <-- fixed
    byline: z.array(z.string()).optional(),
    ingredients: z.array(IngredientsGroup), // <-- fixed Ingredient inside
    suitableForDietIds: z.array(z.string()),
    cuisineIds: z.array(z.string()),
    mealTypeIds: z.array(z.string()),
    celebrationIds: z.array(z.string()),
    utensilsAndApplianceIds: z.array(z.string()),
    techniquesUsedIds: z.array(z.string()),
    difficultyLevel: z.string().optional(),
    serves: z.array(Serves),
    timings: z.array(Timing),
    instructions: z.array(Instruction),
    commerceCtas: z.array(CommerceCTA).optional(), // <-- fixed to require array
    bookCredit: z.string(),
    publishedDate: z.string().optional(),
    firstPublishedDate: z.string().optional(),
    lastModifiedDate: z.string().optional(),
});

export type Recipe = z.infer<typeof Recipe>;