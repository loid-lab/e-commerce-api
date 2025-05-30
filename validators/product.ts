import { z } from "zod";

export const productSchema = z.object({
    title: z.string()
    .trim()
    .max(255, { message: "You have surpassed the 255 character limit" }),
    description: z.string()
    .max(255, { message: "You have surpassed the 255 character limit" }),
    price: z.number()
    .positive({ message: "Price must be greater than zero" })
    .max(10000, { message: "Price is too high" })
    .refine(value => Number.isFinite(value), { message: "Price must be a finite number" })
    .refine(value => {
        return /^\d+(\.\d{1,2})?$/.test(value.toString());
    }, { message: "Price can have up to 2 decimal places" }),
    stock: z.number()
    .int()
    .nonnegative({ message: "Stock cannot be a negative number" }),
    imageUrl: z.string()
    .url({ message: "Image URL must be a valid URL" })
    .max(2048, { message: "ImageURL is too long" })
    .refine(val => val.startsWith("http//") || val.startsWith("https//"), {
        message: "Image URL must start with http:// or https://",
    })
    .optional(),
});

export type ProductSchema = z.infer<typeof productSchema>;