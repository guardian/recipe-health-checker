import {z} from "zod";

export const CheckerOutput = z.object({
    "snapshot": z.string(), //markdown snapshot of the recipe as it was passed to us
    "annotated_text": z.string(),   //snapshot of the recipe with model annotations,
    "model_notes": z.string(),      //optional extra notes
    "annotation_count": z.number(), //how many problems found
    "annotation_summary": z.array(z.string()).nullable().optional(),  //list of the actual problems, corresponds 1:1 with annotation_section
    "annotation_section": z.array(z.string()).nullable().optional(),
    "recipe_id": z.string(),
    "composer_id": z.string(),  //misnamed - actually the CAPI id
    "timestamp": z.string(),
    "input_tokens_used": z.number(),
    "output_tokens_used": z.number(),
    "model_used": z.string(),
})
export type CheckerOutput = z.infer<typeof CheckerOutput>;
