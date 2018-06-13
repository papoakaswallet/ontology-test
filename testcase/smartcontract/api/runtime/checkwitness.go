package runtime

import (
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System.Numerics;

public class A : SmartContract
{
    public static void Main(byte[] input)
    {


        bool b = Runtime.CheckWitness(input);
        if (b) {
            Storage.Put(Storage.CurrentContext, "result", "true");
        }
        Storage.Put(Storage.CurrentContext, "result", "true");
    }
}

code = 53c56b6c766b00527ac4616c766b00c36168184e656f2e52756e74696d652e436865636b5769746e6573736c766b51527ac46c766b51c36c766b52527ac46c766b52c3643f00616168164e656f2e53746f726167652e476574436f6e7465787406726573756c740474727565615272680f4e656f2e53746f726167652e50757461616168164e656f2e53746f726167652e476574436f6e7465787406726573756c740474727565615272680f4e656f2e53746f726167652e50757461616c7566
*/

func TestCheckWitness(ctx *testframework.TestFrameworkContext) bool {
	code := "53c56b6c766b00527ac4616c766b00c36168184e656f2e52756e74696d652e436865636b5769746e6573736c766b51527ac46c766b51c36c766b52527ac46c766b52c3643f00616168164e656f2e53746f726167652e476574436f6e7465787406726573756c740474727565615272680f4e656f2e53746f726167652e50757461616168164e656f2e53746f726167652e476574436f6e7465787406726573756c740474727565615272680f4e656f2e53746f726167652e50757461616c7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestCheckWitness - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestCheckWitness DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	if err != nil {
		ctx.LogError("TestCheckWitness WaitForGenerateBlock error:%s", err)
		return false
	}

	checker, err := ctx.Wallet.NewAccount("", keypair.PK_ECDSA, keypair.P256, signature.SHA256withECDSA, []byte("test"))

	if err != nil {
		ctx.LogError("TestCheckWitness generate account error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		0,
		codeAddress,
		[]interface{}{checker.Address[:]})
	if err != nil {
		ctx.LogError("TestDomainSmartContract InvokeNeoVMSmartContract error: %s", err)
	}

	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	res, err := ctx.Ont.Rpc.GetStorage(codeAddress, []byte("result"))
	if err != nil {
		ctx.LogError("TestDomainSmartContract GetStorageItem key:hello error: %s", err)
		return false
	}

	err = ctx.AssertToString(string(res), "true")
	if err != nil {
		ctx.LogError("TestDomainSmartContract AssertToString error: %s", err)
		return false
	}
	return true
}
