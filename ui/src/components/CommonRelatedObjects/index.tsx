/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { memo, FC } from 'react';
import { Card, ListGroup } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { Icon } from '@/components';
import { useSimilarArticle } from '@/services';
import { pathFactory } from '@/router/pathFactory';

import type * as Type from '@/common/interface';


interface Props {
  // id: string;
  head_title:string;
  // data:Type.ListResult;
  // isLoading :boolean;
  // error: Error;
  gen_link_func: (item:any)=> string;
  get_data_func: () => any;

}
const Index: FC<Props> = ({ 
  // id,
  head_title,
  get_data_func,
  // data,
  // isLoading,
  // error,
  gen_link_func,
  

 }) => {
  // const { t } = useTranslation('translation', {
  //   keyPrefix: 'related_article',
  // });

  // const { data, isLoading } = useSimilarArticle({
  //   article_id : id,
  //   page_size: 5,
  // });
  const { data, isLoading } = get_data_func()
  if (isLoading) {
    return null;
  }

  return (
    <Card>
      <Card.Header>{head_title}</Card.Header>
      <ListGroup variant="flush">
        {data?.list?.map((item) => {
          return (
            <ListGroup.Item
              action
              key={item.id}
              as={Link}
              to={gen_link_func(item)}>
              <div className="link-dark">{item.title}</div>
              
            </ListGroup.Item>
          );
        })}
      </ListGroup>
    </Card>
  );
};

export default memo(Index);
