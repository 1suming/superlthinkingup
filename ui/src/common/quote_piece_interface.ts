import { UserInfoBase,ImgCodeReq,Tag,Paging } from './interface';


// export interface QuotePieceDetailRes {
//     id: string;
//     title: string;
//     content: string;
//     html: string;
//     tags: any[];
//     view_count: number;
//     unique_view_count?: number;
  
//     create_time: string;
//     update_time: string;
//     user_info: UserInfoBase;
  
//     [prop: string]: any;
//   }
  export interface QuotePieceDetailRes {
    id: string;
    title: string;
    content: string;
    html: string;
    tags: any[];
    view_count: number;
    unique_view_count?: number;
    answer_count: number;
    favorites_count: number;
    follow_counts: 0;
    accepted_answer_id: string;
    last_answer_id: string;
    create_time: string;
    update_time: string;
    user_info: UserInfoBase;
    answered: boolean;
    collected: boolean;
    answer_ids: string[];
  
   
    [prop: string]: any;

    content_format: number;
  }


  export interface QuotePieceParams extends ImgCodeReq {
    title: string;
    url_title?: string;
    content: string;
    tags: Tag[];

    content_format: number;
  }

  export interface QuotePieceOperationReq {
    id: string;
    operation: 'pin' | 'unpin' | 'hide' | 'show';
  }
  
// export interface QuerySiteInfoKeyValReq  {
//     key: string; 
// }
// export interface QuerySiteInfoKeyValResp   {
//     content: string;
// }
export type QuotePieceOrderBy =
  | 'recommend'
  | 'newest'
  | 'active'
  | 'hot'
  | 'score'
  | 'unanswered';
export interface QueryQuotePiecesReq extends Paging {
    order: QuotePieceOrderBy;
    tag?: string;
    in_days?: number;

    tag_id?:string;
}