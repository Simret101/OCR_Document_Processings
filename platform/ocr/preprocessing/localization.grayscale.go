package preprocessing

import "gocv.io/x/gocv"

 func (p *Storage) Gray(src gocv.Mat) (gocv.Mat, error) {                   
       if src.Empty() {                                                         
         return gocv.NewMat(), nil                                              
       }                                                                        
       if src.Channels() == 1 {                                                 
         return src.Clone(), nil                                                
       }                                                                        
                                                                                
       dst := gocv.NewMat()                                                     
       gocv.CvtColor(src, &dst, gocv.ColorBGRToGray)                            
       return dst, nil                                                          
     }  