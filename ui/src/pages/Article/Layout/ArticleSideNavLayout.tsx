import { FC, memo } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import { Outlet } from 'react-router-dom';

import  ArticleSideNav from '@/components/Article/SideNav/ArticleSideNav';

import '@/css/articleSideLayout.scss';

const Index: FC = () => {

    const  myStyle={
        'height':"100vh",
    }
    //ArticleSideNav 不算侧边栏了

  return (
    <Container className="d-flex flex-column flex-fill">
        
           
                <ArticleSideNav />   
            
         
      <Row className="flex-fill">
        <Col xl={1} ></Col>
        <Col xl={10} >
          <Outlet />
        </Col>
        <Col xl={1} ></Col>
      </Row>
    </Container>
  );
};

export default memo(Index);
