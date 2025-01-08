import { FC, memo,useEffect } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import { Outlet,useLocation } from 'react-router-dom';
 
import  QuoteSideNav from '@/components/Quote/SideNav/QuoteSideNav';

import '@/css/quote_module.scss';

const Index: FC = () => {
   
  const location = useLocation();

  //动态修改classlist
  useEffect(() => {
      document.body.classList.add("body-main-content-quote");
      return ()=>{
        document.body.classList.remove('body-main-content-quote');
      };
  }, []);
  
    const  myStyle={
        'height':"100vh",
    }
    //ArticleSideNav 不算侧边栏了

  return (
    <Container className="d-flex flex-column flex-fill">
        
           
                <QuoteSideNav />   
            
         
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
