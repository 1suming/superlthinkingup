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

import { FC ,useState,useEffect} from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {Row, Col, Nav } from 'react-bootstrap';
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';
import NavDropdown from 'react-bootstrap/NavDropdown';
import classnames from 'classnames';

import { loggedUserInfoStore, sideNavStore,myGlobalInfoStore } from '@/stores';
import { Icon } from '@/components';
import './index.scss';






const Index: FC = () => {
  const { t } = useTranslation();
  const { pathname } = useLocation();
  const { user: userInfo } = loggedUserInfoStore();
  const { visible, can_revision, revision } = sideNavStore();

  const {isSideNavSticky,sideNavStickyTop}= myGlobalInfoStore()

//   const [scrolled,setScrolled]=useState(false);
//   const handleScroll=() => {
//       const offset=window.scrollY;
      
//       if(offset > 100 ){
//         setScrolled(true);
//         console.log("offset scorll>>")
//       }
//       else{
//         setScrolled(false);
//         console.log("offset scorll false")
//       }
//     }
  
//     useEffect(() => {
//       window.addEventListener('scroll',handleScroll)
//     });
//     let scrolledCls="";
//     if(scrolled){
//       scrolledCls="stick-top";
//     }else{
//       scrolledCls="";
//     }
    const  sideNavStickTopStyle={
        'position':"sticky",
        'top': sideNavStickyTop,
        'zIndex':"1019", //最顶部的top是1020，要比它小，不然会把top的box-shadow挡住

        'background':'#fff',
 
    } as React.CSSProperties ; // 'background':'#fff',
    const emptyStyle={} as React.CSSProperties;

    // style={ isSideNavSticky?sideNavStickTopStyle:emptyStyle }

    console.log("二级菜单:",sideNavStickyTop);
// className="`stick-top bg-body-tertiary`
// position:sticky;失效原因：https://www.cnblogs.com/coco1s/p/14180476.html
return (
    <>
    </>
);

//   return (
//     <Row className="flex-fill" style={ sideNavStickTopStyle }>
//          <Col xl={1} ></Col>
//          <Col xl={10} >
//         <nav id="second-article-sideNav" className="nav"   >
//             <NavLink to="/questions" className="nav-link">全部 </NavLink>
//             <a className="nav-link" href="#">思维模型</a>
//             <a className="nav-link" href="#">个人提升</a>
//         </nav>
//         </Col>
//         <Col xl={1} ></Col>
//     </Row>
//   );
};

export default Index;
