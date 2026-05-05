import type {contentservicev1_Category} from "@/api/generated/app/service/v1";

export interface Props {
  categories: contentservicev1_Category[]
  level?: number
}
