import React, { useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
 
// // 自定义 Hook 修改 body 样式
// const useBodyBackground = () => {
//   const background = useBackgroundStore((state) => state.background); // 获取背景颜色

//   useEffect(() => {
//     document.body.style.backgroundColor = background;

//     return () => {
//       document.body.style.backgroundColor = "";
//     };
//   }, [background]);
// };