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

import { FC, memo, useState, useEffect,useRef} from 'react';
import {
  Navbar,
  Container,
  Nav,
  Form,
  FormControl,
  Col,
  Dropdown,
  NavDropdown,
} from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import {
  useSearchParams,
  Link,
  useNavigate,
  useLocation,
  useMatch,
  NavLink,
} from 'react-router-dom';

import classnames from 'classnames';



import { userCenter, floppyNavigation } from '@/utils';
import {
  loggedUserInfoStore,
  siteInfoStore,
  brandingStore,
  loginSettingStore,
  themeSettingStore,
  sideNavStore,
} from '@/stores';
import { logout, useQueryNotificationStatus } from '@/services';
import { Icon } from '@/components';

import NavItems from './components/NavItems';

import './index.scss';

interface  HeaderProps {
    siteHeadNavRef: React.RefObject<HTMLElement>;
}
const Header: FC<HeaderProps> = ( {siteHeadNavRef}) => {
  const navigate = useNavigate();
  const location = useLocation();
  const [urlSearch] = useSearchParams();
  const q = urlSearch.get('q');
  const { user, clear: clearUserStore } = loggedUserInfoStore();
  const { t } = useTranslation();
  const [searchStr, setSearch] = useState('');
  const siteInfo = siteInfoStore((state) => state.siteInfo);
  const brandingInfo = brandingStore((state) => state.branding);
  const loginSetting = loginSettingStore((state) => state.login);
  const { updateReview, updateVisible } = sideNavStore();
  const { data: redDot } = useQueryNotificationStatus();
  /**
   * Automatically append `tag` information when creating a question
   */
  const tagMatch = useMatch('/tags/:slugName');
  let askUrl = '/questions/ask';
  if (tagMatch && tagMatch.params.slugName) {
    askUrl = `${askUrl}?tags=${encodeURIComponent(tagMatch.params.slugName)}`;
  }

  useEffect(() => {
    updateReview({
      can_revision: Boolean(redDot?.can_revision),
      revision: Number(redDot?.revision),
    });
  }, [redDot]);

  const handleInput = (val) => {
    setSearch(val);
  };
  const handleSearch = (evt) => {
    evt.preventDefault();
    if (!searchStr) {
      return;
    }
    const searchUrl = `/search?q=${encodeURIComponent(searchStr)}`;
    navigate(searchUrl);
  };

  const handleLogout = async (evt) => {
    evt.preventDefault();
    await logout();
    clearUserStore();
    window.location.replace(window.location.href);
  };

  useEffect(() => {
    if (q && location.pathname === '/search') {
      handleInput(q);
    }
  }, [q]);

  useEffect(() => {
    const collapse = document.querySelector('#navBarContent');
    if (collapse && collapse.classList.contains('show')) {
      const toggle = document.querySelector('#navBarToggle') as HTMLElement;
      if (toggle) {
        toggle?.click();
      }
    }

    // clear search input when navigate to other page
    if (location.pathname !== '/search' && searchStr) {
      setSearch('');
    }
    
  }, [location.pathname]);

  let navbarStyle = 'theme-colored';
  const { theme, theme_config } = themeSettingStore((_) => _);
  if (theme_config?.[theme]?.navbar_style) {
    navbarStyle = `theme-${theme_config[theme].navbar_style}`;
  }
  console.log('@navbarStyle:',navbarStyle)
//   const ax=useRef(null);

//   useEffect(() => {
//     const ex= ax.current;
//     console.log(ex);
//   });

    const { pathname } = useLocation();
    console.log("pathname:",pathname);

  const goToHref = function(event,href){
    event.preventDefault();
   // console.log("event.target.href:",event.target.href)//,event.target.href会加上http域名，不符合要求
    navigate(href);
  };
  return (
    <Navbar ref={siteHeadNavRef}
      variant={navbarStyle === 'theme-colored' ? 'dark' : ''}
      expand="lg"
      className={classnames('sticky-top', navbarStyle)}
      id="header">
      <Container className="d-flex align-items-center">
        <Navbar.Toggle
          aria-controls="navBarContent"
          className="answer-navBar me-2"
          id="navBarToggle"
          onClick={() => {
            updateVisible();
          }}
        />

        <div className="d-flex justify-content-between align-items-center nav-grow flex-nowrap">
          <Navbar.Brand to="/" as={Link} className="lh-1 me-0 me-sm-5 p-0">
            {brandingInfo.logo ? (
              <>
                <img
                  className="d-none d-lg-block logo me-0"
                  src={brandingInfo.logo}
                  alt={siteInfo.name}
                />

                <img
                  className="lg-none logo me-0"
                  src={brandingInfo.mobile_logo || brandingInfo.logo}
                  alt={siteInfo.name}
                />
              </>
            ) : (
              <span>{siteInfo.name}</span>
            )}
          </Navbar.Brand>

          {/* mobile nav */}
          <div className="d-flex lg-none align-items-center flex-lg-nowrap">
            {user?.username ? (
              <NavItems
                redDot={redDot}
                userInfo={user}
                logOut={(e) => handleLogout(e)}
              />
            ) : (
              <>
                <Link
                  className={classnames('me-2 btn btn-link', {
                    'link-light': navbarStyle === 'theme-colored',
                    'link-primary': navbarStyle !== 'theme-colored',
                  })}
                  onClick={() => floppyNavigation.storageLoginRedirect()}
                  to={userCenter.getLoginUrl()}>
                  {t('btns.login')}
                </Link>
                {loginSetting.allow_new_registrations && (
                  <Link
                    className={classnames(
                      'btn',
                      navbarStyle === 'theme-colored'
                        ? 'btn-light'
                        : 'btn-primary',
                    )}
                    to={userCenter.getSignUpUrl()}>
                    {t('btns.signup')}
                  </Link>
                )}
              </>
            )}
          </div>
        </div>




        <Navbar.Collapse id="navBarContent" className="me-auto">
          <hr className="hr lg-none mb-3" style={{ marginTop: '12px' }} />
          <Col lg={5} className="ps-0 nav-left-top-bar">
           <Nav className="nav-top">
                 
 
                <NavLink
                    to="/"
                    className={() =>
                    pathname === '/' ? 'nav-link active' : 'nav-link'
                    }>
                    <span>首页</span>
                </NavLink>

                <NavLink
                    to="/articles"
                    className={() =>
                    pathname === '/articles' ? 'nav-link active' : 'nav-link'
                    }>
                    <span>热文</span>
                </NavLink>
                <NavLink
                    to="/questions"
                    className={() =>
                    pathname === '/questions' ? 'nav-link active' : 'nav-link'
                    }>
                    <span>问答</span>
                </NavLink>
                <NavLink
                    to="/quotes"
                    className={() =>
                    pathname === '/quotes' ? 'nav-link active' : 'nav-link'
                    }>
                    <span>金句控</span>
                </NavLink>
               <NavDropdown title="人工智能AI" id="navbarScrollingDropdown1">
                
                <NavDropdown.Item href="/ai" onClick={(evt) => { goToHref(evt,"/ai");}}>总览</NavDropdown.Item>
                <NavDropdown.Divider />
                           </NavDropdown>

             <NavDropdown title="更多" id="navbarScrollingDropdown">
                
                <NavDropdown.Item href="/siteinfo/about" onClick={(evt) => { goToHref(evt,"/siteinfo/about");}}>超维社是什么?</NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item href="/siteinfo/about" onClick={(evt) => { goToHref(evt,"/siteinfo/about");}}> 关于超维社 </NavDropdown.Item>
                <NavDropdown.Item href="/siteinfo/contact" onClick={(evt) => { goToHref(evt,"/siteinfo/contact");}}> 联系我们 </NavDropdown.Item>
                <NavDropdown.Item href="/siteinfo/disclaim" onClick={(evt) => { goToHref(evt,"/siteinfo/disclaim");}} > 免责声明 </NavDropdown.Item>
            </NavDropdown>
 
           

            </Nav>
          </Col>




          <Col lg={3} className="ps-0">
            <Form
              action="/search"
              className="w-100 maxw-400 position-relative"
              onSubmit={handleSearch}>
              <div className="search-wrap" onClick={handleSearch}>
                <Icon name="search" className="search-icon" />
              </div>
              <FormControl
                type="search"
                placeholder={t('header.search.placeholder')}
                className="placeholder-search"
                value={searchStr}
                name="q"
                onChange={(e) => handleInput(e.target.value)}
              />
            </Form>
          </Col>

          <Nav.Item className="lg-none mt-3 pb-1">
            <Link
              to={askUrl}
              className="text-capitalize text-nowrap btn btn-light">
              {t('btns.add_question')}
            </Link>
          </Nav.Item>
          {/* pc nav */}
          <Col
            lg={4}
            className="d-none d-lg-flex justify-content-start justify-content-sm-end">
            {user?.username ? (
              <Nav className="d-flex align-items-center flex-lg-nowrap">
                <Nav.Item className="me-3">
                  <Link
                    to={askUrl}
                    className={classnames('text-capitalize text-nowrap btn', {
                      'btn-light': navbarStyle !== 'theme-light',
                      'btn-primary': navbarStyle === 'theme-light',
                    })}>
                    {t('btns.add_question')}
                  </Link>
                </Nav.Item>

                <NavItems
                  redDot={redDot}
                  userInfo={user}
                  logOut={handleLogout}
                />
              </Nav>
            ) : (
              <>
                <Link
                  className={classnames('me-2 btn btn-link', {
                    'link-light': navbarStyle === 'theme-colored',
                    'link-primary': navbarStyle !== 'theme-colored',
                  })}
                  onClick={() => floppyNavigation.storageLoginRedirect()}
                  to={userCenter.getLoginUrl()}>
                  {t('btns.login')}
                </Link>
                {loginSetting.allow_new_registrations && (
                  <Link
                    className={classnames(
                      'btn',
                      navbarStyle === 'theme-colored'
                        ? 'btn-light'
                        : 'btn-primary',
                    )}
                    to={userCenter.getSignUpUrl()}>
                    {t('btns.signup')}
                  </Link>
                )}
              </>
            )}
          </Col>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default memo(Header);
