export { };
declare global {
    enum Recurring {
        No = "no",
        Daily = "daily",
        Weekly = "weekly",
        Monthly = "monthly",
        Yearly = "yearly",
    }

    type newItem = {
        name: string,
        price: number,
        saved: number,
        recurring: Recurring,
        section_id: number,
    }
}
