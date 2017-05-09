<?php

/**
 * btscript  game1800001 /tmp/TestMatch.php 20001 Match game1800001 pirate_warship_36 1800001 10
 */
class TestMatch extends BaseScript
{

    /* (non-PHPdoc)
     * @see BaseScript::executeScript()
     */
    protected function executeScript($arrOption)
    {

        $uid = 0;
        if (isset ( $arrOption [0] ))
        {
            $uid = intval ( $arrOption [0] );
        }
        else
        {
            exit ( "usage: uid method val1 val2 \n" );
        }
        
        $ret=array();
        $proxy = new PHPProxy ( 'match' );
        if ($arrOption [1]=="Sum"){
            $request=array(
                "userid"=>$uid,
                "num1"=>$arrOption[2],
                "num2"=>$arrOption[3],
            );
            $ret=$proxy->Sum(Util::amfEncode($request));
        }elseif ($arrOption [1]=="Concat"){
            $request= array(
                "userid"=>$uid,
                "str1"=>$arrOption[2],
                "str2"=>$arrOption[3],
            );
            $ret=$proxy->Concat(Util::amfEncode($request));
        } elseif ($arrOption [1]=="Match"){
            $request= array(
                "gameGroup"=>$arrOption[2],
                "gameDb"=>$arrOption[3],
                "serverId"=>$arrOption[4],
                "userid"=>$uid,
                "score"=>$arrOption[5],
            );
            $ret=$proxy->Match(Util::amfEncode($request));
        }
        var_dump(Util::amfDecode($ret));
    }
    
 }

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */