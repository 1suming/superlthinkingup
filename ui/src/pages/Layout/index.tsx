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

import { FC, memo, useEffect ,useRef,useState} from 'react';
import { Outlet, useLocation, ScrollRestoration } from 'react-router-dom';
import { HelmetProvider } from 'react-helmet-async';

import { SWRConfig } from 'swr';
//import { throttle, debounce } from 'lodash';
import debounce from 'lodash/debounce'; //@cws 按需加载
import throttle from 'lodash/throttle';

import { toastStore, loginToContinueStore, errorCodeStore,myGlobalInfoStore } from '@/stores';
import {
  Header,
  Footer,
  Toast,
  Customize,
  CustomizeTheme,
  PageTags,
  HttpErrorContent,
} from '@/components';
import { LoginToContinueModal, BadgeModal } from '@/components/Modal';
import { changeTheme } from '@/utils';
import { useQueryNotificationStatus } from '@/services';

 


const Layout: FC = () => {
  const location = useLocation();
  const { msg: toastMsg, variant, clear: toastClear } = toastStore();
  const closeToast = () => {
    toastClear();
  };
  const { code: httpStatusCode, reset: httpStatusReset } = errorCodeStore();
  const { show: showLoginToContinueModal } = loginToContinueStore();
  const { data: notificationData } = useQueryNotificationStatus();



  useEffect(() => {
    httpStatusReset();
  }, [location]);

  useEffect(() => {
    const systemThemeQuery = window.matchMedia('(prefers-color-scheme: dark)');
    function handleSystemThemeChange(event) {
      if (event.matches) {
        changeTheme('dark');
      } else {
        changeTheme('light');
      }
    }

    systemThemeQuery.addListener(handleSystemThemeChange);

    return () => {
      systemThemeQuery.removeListener(handleSystemThemeChange);
    };
  }, []);

   //@cws
 const siteHeadNavRef=useRef< HTMLElement>(null);

//   const [siteHeadNavHeight,setSiteHeadNavHeight]= useState(0);
  let siteHeadNavHeight=0; //setState每次调用都会重新渲染，不需要用setState

//   const [isSideNavSticky,setIsSideNavSticky]=useState(false);
  const {isSideNavSticky,setIsSideNavSticky,setSideNavStickyTop}= myGlobalInfoStore()
  const handleScroll=() => {
    const offset=window.scrollY;
    if(siteHeadNavHeight!==undefined){
        if(offset > siteHeadNavHeight ){
        //   isSideNavSticky=true;
         setIsSideNavSticky(true);
          setSideNavStickyTop(siteHeadNavHeight);
          console.log("offset > sitehadNavHeight scorll true:"+"offset:"+offset+",siteHeadNavHeight"+siteHeadNavHeight + " isSideNavSticky,"+ isSideNavSticky);
        }
        else{
            setIsSideNavSticky(false);
            console.log("offset > sitehadNavHeight scorll  false:",isSideNavSticky);
        }
    }
 
  }
    // 使用节流
  const throttledScrollHandler = throttle(handleScroll, 200);
    // // 使用防抖
    // const debouncedScrollHandler = debounce(handleScroll, 200);
  

    useEffect(() => {
        if (siteHeadNavRef.current) {
            const { height } = siteHeadNavRef.current.getBoundingClientRect();
            if (height !== undefined && height > 0) {
                // setSiteHeadNavHeight(height);
                siteHeadNavHeight=height;
                console.log("siteHeadNavHeight>>", height);
                //总是sticky
                setIsSideNavSticky(true);
                setSideNavStickyTop(siteHeadNavHeight);
            }
            else { console.log("siteHeadNavHeight is undefined"); }
        }   
  
        // window.addEventListener('scroll',throttledScrollHandler)
        // return () => {
        //     window.removeEventListener('scroll', throttledScrollHandler);
        // };
    }, []);  
 
//   let scrolledCls="";
//   if(scrolled){
//     scrolledCls="stick-top";
//   }else{
//     scrolledCls="";
//   }


  return (
    <HelmetProvider>
      <PageTags />
      <CustomizeTheme />
      <SWRConfig
        value={{
          revalidateOnFocus: false,
        }}>
        <Header siteHeadNavRef={siteHeadNavRef}/>
        {/* eslint-disable-next-line jsx-a11y/click-events-have-key-events */}
        <div className="position-relative page-wrap d-flex flex-column flex-fill">
          {httpStatusCode ? (
            <HttpErrorContent httpCode={httpStatusCode} />
          ) : (
            <Outlet />
          )}
        </div>
        <Toast msg={toastMsg} variant={variant} onClose={closeToast} />
        <Footer />
        <Customize />
        <LoginToContinueModal visible={showLoginToContinueModal} />
        <BadgeModal
          badge={notificationData?.badge_award}
          visible={Boolean(notificationData?.badge_award)}
        />
        <ScrollRestoration />
      </SWRConfig>
    </HelmetProvider>
  );
};

export default memo(Layout);
