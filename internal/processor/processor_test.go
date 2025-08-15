package processor_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
	"github.com/atmxlab/atmcfg/internal/linker"
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/processor"
	"github.com/atmxlab/atmcfg/internal/processor/mocks"
	"github.com/atmxlab/atmcfg/internal/test/testast"
	"github.com/atmxlab/atmcfg/internal/test/testlinkedast"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TestProcessSuite struct {
	suite.Suite
	osMock     *mocks.MockOS
	lexerMock  *mocks.MockLexer
	parserMock *mocks.MockParser
	linkerMock *mocks.MockLinker

	processer *processor.Processor
}

func (s *TestProcessSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.osMock = mocks.NewMockOS(ctrl)
	s.lexerMock = mocks.NewMockLexer(ctrl)
	s.parserMock = mocks.NewMockParser(ctrl)
	s.linkerMock = mocks.NewMockLinker(ctrl)

	s.processer = processor.NewProcessor(
		s.lexerMock,
		s.parserMock,
		s.linkerMock,
		s.osMock,
	)
}

func TestProcess(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TestProcessSuite))
}

func (s *TestProcessSuite) TestHappyPath() {
	s.T().Parallel()

	// Arrange
	var (
		configPath    = "./dir1/dir2/config.atmc"
		absConfigPath = "/home/user/dir1/dir2/config.atmc"
		fileContent   = "{a: 123, b: 11}"
	)

	s.osMock.EXPECT().AbsPath(configPath, ".").Return(absConfigPath, nil)
	s.osMock.EXPECT().ReadFile(absConfigPath).Return([]byte(fileContent), nil)

	s.lexerMock.EXPECT().Tokenize(gomock.Any()).Return(&tokenmover.TokenMover{}, nil)

	s.parserMock.EXPECT().Parse(gomock.Any()).Return(testast.NewAstBuilder().Build(), nil)

	expectedMainAstWithPath := ast.NewWithPath(
		testast.NewAstBuilder().Build(),
		absConfigPath,
		map[string]string{},
	)
	s.linkerMock.EXPECT().Link(gomock.Any()).DoAndReturn(func(param linker.LinkParam) (linkedast.Ast, error) {
		expectedParam := linker.LinkParam{
			MainAst: expectedMainAstWithPath,
			ASTByPath: map[string]ast.WithPath{
				absConfigPath: expectedMainAstWithPath,
			},
			Env: nil,
		}

		// Assert
		require.Equal(s.T(), expectedParam.MainAst, param.MainAst)
		require.Equal(s.T(), expectedParam.ASTByPath, param.ASTByPath)
		require.Equal(s.T(), expectedParam.Env, param.Env)

		return testlinkedast.NewAstBuilder().Build(), nil
	})

	// Act
	err := s.processer.Process(configPath)

	// Assert
	require.NoError(s.T(), err)
}
